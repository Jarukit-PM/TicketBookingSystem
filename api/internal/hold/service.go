package hold

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"time"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/booking"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/catalog"
)

// Service manages Redis-backed seat holds for showtimes.
type Service struct {
	showtimes catalog.ShowtimeRepository
	screens   catalog.ScreenRepository
	cinemas   catalog.CinemaRepository
	bookings  booking.Repository
	redis     *redis.Client
	ttl       time.Duration
	now       func() time.Time
}

// Option configures a hold Service.
type Option func(*Service)

// WithClock overrides the clock used for TTL and showtime cutoff checks.
func WithClock(now func() time.Time) Option {
	return func(s *Service) { s.now = now }
}

// NewService returns a hold service with the default 5-minute TTL.
func NewService(
	showtimes catalog.ShowtimeRepository,
	screens catalog.ScreenRepository,
	cinemas catalog.CinemaRepository,
	bookings booking.Repository,
	rdb *redis.Client,
	opts ...Option,
) *Service {
	s := &Service{
		showtimes: showtimes,
		screens:   screens,
		cinemas:   cinemas,
		bookings:  bookings,
		redis:     rdb,
		ttl:       DefaultTTL,
		now:       time.Now,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

type showtimeCtx struct {
	showtimeID string
	showtime   *catalog.Showtime
	screen     *catalog.Screen
	layout     map[string]catalog.LayoutSeat
	sold       map[string]struct{}
}

// AddSeats holds additional seats for a user, refreshing TTL on all of their holds.
func (s *Service) AddSeats(ctx context.Context, userID, showtimeID string, seatIDs []string) (Result, error) {
	if s.redis == nil {
		return Result{}, fmt.Errorf("redis client is nil")
	}

	stx, err := s.loadShowtimeCtx(ctx, showtimeID)
	if err != nil {
		return Result{}, err
	}
	if err := s.ensureShowtimeOpen(ctx, stx); err != nil {
		return Result{}, err
	}

	seatIDs = uniqueSorted(seatIDs)
	if len(seatIDs) == 0 {
		return s.currentResult(ctx, userID, showtimeID)
	}

	current, err := s.listUserHolds(ctx, userID, showtimeID)
	if err != nil {
		return Result{}, err
	}
	currentSet := toSet(current)

	for _, seatID := range seatIDs {
		if err := s.validateSeatForAdd(ctx, stx, userID, showtimeID, seatID); err != nil {
			return Result{}, err
		}
	}

	merged := uniqueSorted(append(current, seatIDs...))
	if len(merged) > MaxSeatsPerHold {
		return Result{}, ErrSeatLimitExceeded
	}

	heldAt := s.now().UTC()
	for _, seatID := range seatIDs {
		if _, ok := currentSet[seatID]; ok {
			continue
		}
		payload, err := json.Marshal(Record{UserID: userID, HeldAt: heldAt})
		if err != nil {
			return Result{}, fmt.Errorf("marshal hold: %w", err)
		}
		ok, err := s.redis.SetNX(ctx, SeatKey(showtimeID, seatID), payload, s.ttl).Result()
		if err != nil {
			return Result{}, fmt.Errorf("set hold: %w", err)
		}
		if !ok {
			owner, ownerErr := s.holdOwner(ctx, showtimeID, seatID)
			if ownerErr != nil {
				return Result{}, ownerErr
			}
			if owner != userID {
				return Result{}, ErrSeatHeldByOther
			}
		}
	}

	if len(seatIDs) > 0 {
		members := make([]any, len(seatIDs))
		for i, seatID := range seatIDs {
			members[i] = seatID
		}
		if err := s.redis.SAdd(ctx, UserHoldsKey(userID, showtimeID), members...).Err(); err != nil {
			return Result{}, fmt.Errorf("track user holds: %w", err)
		}
	}

	expiresAt, err := s.refreshUserHoldsTTL(ctx, userID, showtimeID)
	if err != nil {
		return Result{}, err
	}

	added := make([]string, 0, len(seatIDs))
	for _, seatID := range seatIDs {
		if _, wasHeld := currentSet[seatID]; !wasHeld {
			added = append(added, seatID)
		}
	}

	return Result{Holds: merged, ExpiresAt: &expiresAt, Added: added}, nil
}

// RemoveSeats releases seats without refreshing TTL on remaining holds.
// An empty seatIDs slice releases all holds for the showtime.
func (s *Service) RemoveSeats(ctx context.Context, userID, showtimeID string, seatIDs []string) (Result, error) {
	if s.redis == nil {
		return Result{}, fmt.Errorf("redis client is nil")
	}

	current, err := s.listUserHolds(ctx, userID, showtimeID)
	if err != nil {
		return Result{}, err
	}
	if len(current) == 0 {
		return Result{Holds: []string{}}, nil
	}

	toRemove := uniqueSorted(seatIDs)
	if len(toRemove) == 0 {
		toRemove = current
	}

	currentSet := toSet(current)
	for _, seatID := range toRemove {
		if _, ok := currentSet[seatID]; !ok {
			return Result{}, ErrSeatNotHeld
		}
	}

	userKey := UserHoldsKey(userID, showtimeID)
	pipe := s.redis.Pipeline()
	for _, seatID := range toRemove {
		pipe.Del(ctx, SeatKey(showtimeID, seatID))
		pipe.SRem(ctx, userKey, seatID)
	}
	if _, err := pipe.Exec(ctx); err != nil {
		return Result{}, fmt.Errorf("release holds: %w", err)
	}

	removeSet := toSet(toRemove)
	remaining := make([]string, 0, len(current)-len(toRemove))
	for _, seatID := range current {
		if _, drop := removeSet[seatID]; !drop {
			remaining = append(remaining, seatID)
		}
	}
	remaining = uniqueSorted(remaining)

	if len(remaining) == 0 {
		if err := s.redis.Del(ctx, userKey).Err(); err != nil {
			return Result{}, fmt.Errorf("clear user holds: %w", err)
		}
		return Result{Holds: []string{}, Released: toRemove}, nil
	}

	ttl, err := s.redis.TTL(ctx, userKey).Result()
	if err != nil {
		return Result{}, fmt.Errorf("read hold ttl: %w", err)
	}

	var expiresAt *time.Time
	if ttl > 0 {
		t := s.now().Add(ttl)
		expiresAt = &t
	}

	return Result{Holds: remaining, ExpiresAt: expiresAt, Released: toRemove}, nil
}

// UserHolds returns the sorted seat IDs held by a user on a showtime.
func (s *Service) UserHolds(ctx context.Context, userID, showtimeID string) ([]string, error) {
	return s.listUserHolds(ctx, userID, showtimeID)
}

// ClearUserHolds releases all holds for a user on a showtime.
func (s *Service) ClearUserHolds(ctx context.Context, userID, showtimeID string) error {
	_, err := s.RemoveSeats(ctx, userID, showtimeID, nil)
	return err
}

func (s *Service) currentResult(ctx context.Context, userID, showtimeID string) (Result, error) {
	holds, err := s.listUserHolds(ctx, userID, showtimeID)
	if err != nil {
		return Result{}, err
	}
	if len(holds) == 0 {
		return Result{Holds: []string{}}, nil
	}

	ttl, err := s.redis.TTL(ctx, UserHoldsKey(userID, showtimeID)).Result()
	if err != nil {
		return Result{}, fmt.Errorf("read hold ttl: %w", err)
	}
	if ttl <= 0 {
		return Result{Holds: holds}, nil
	}
	expiresAt := s.now().Add(ttl)
	return Result{Holds: holds, ExpiresAt: &expiresAt}, nil
}

func (s *Service) loadShowtimeCtx(ctx context.Context, showtimeID string) (*showtimeCtx, error) {
	oid, err := primitive.ObjectIDFromHex(showtimeID)
	if err != nil {
		return nil, ErrShowtimeNotFound
	}

	showtime, err := s.showtimes.FindShowtimeByID(ctx, oid)
	if err != nil {
		return nil, fmt.Errorf("find showtime: %w", err)
	}
	if showtime == nil {
		return nil, ErrShowtimeNotFound
	}

	screen, err := s.screens.FindScreenByID(ctx, showtime.ScreenID)
	if err != nil {
		return nil, fmt.Errorf("find screen: %w", err)
	}
	if screen == nil {
		return nil, ErrScreenNotFound
	}

	bookings, err := s.bookings.ListConfirmedByShowtime(ctx, oid)
	if err != nil {
		return nil, fmt.Errorf("list bookings: %w", err)
	}

	layout := make(map[string]catalog.LayoutSeat, len(screen.Layout.Seats))
	for _, seat := range screen.Layout.Seats {
		layout[seat.SeatID] = seat
	}

	sold := make(map[string]struct{})
	for _, b := range bookings {
		for _, seatID := range b.Seats {
			sold[seatID] = struct{}{}
		}
	}

	return &showtimeCtx{
		showtimeID: showtimeID,
		showtime:   showtime,
		screen:     screen,
		layout:     layout,
		sold:       sold,
	}, nil
}

func (s *Service) ensureShowtimeOpen(ctx context.Context, stx *showtimeCtx) error {
	cinema, err := s.cinemas.FindCinemaByID(ctx, stx.screen.CinemaID)
	if err != nil {
		return fmt.Errorf("find cinema: %w", err)
	}
	if cinema == nil {
		return ErrCinemaNotFound
	}

	loc, err := time.LoadLocation(cinema.Timezone)
	if err != nil {
		loc = time.UTC
	}

	now := s.now().In(loc)
	startsAt := stx.showtime.StartsAt.In(loc)
	if !startsAt.After(now) {
		return ErrShowtimeStarted
	}
	return nil
}

func (s *Service) validateSeatForAdd(
	ctx context.Context,
	stx *showtimeCtx,
	userID, showtimeID, seatID string,
) error {
	seat, ok := stx.layout[seatID]
	if !ok {
		return ErrSeatNotFound
	}
	if seat.Type == catalog.SeatTypeBlocked {
		return ErrSeatBlocked
	}
	if _, sold := stx.sold[seatID]; sold {
		return ErrSeatSold
	}

	owner, err := s.holdOwner(ctx, showtimeID, seatID)
	if err != nil {
		return err
	}
	if owner != "" && owner != userID {
		return ErrSeatHeldByOther
	}
	return nil
}

func (s *Service) holdOwner(ctx context.Context, showtimeID, seatID string) (string, error) {
	raw, err := s.redis.Get(ctx, SeatKey(showtimeID, seatID)).Result()
	if err == redis.Nil {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("get hold: %w", err)
	}

	var rec Record
	if err := json.Unmarshal([]byte(raw), &rec); err != nil {
		return "", fmt.Errorf("decode hold: %w", err)
	}
	return rec.UserID, nil
}

func (s *Service) listUserHolds(ctx context.Context, userID, showtimeID string) ([]string, error) {
	members, err := s.redis.SMembers(ctx, UserHoldsKey(userID, showtimeID)).Result()
	if err != nil {
		return nil, fmt.Errorf("list user holds: %w", err)
	}
	return uniqueSorted(members), nil
}

func (s *Service) refreshUserHoldsTTL(ctx context.Context, userID, showtimeID string) (time.Time, error) {
	holds, err := s.listUserHolds(ctx, userID, showtimeID)
	if err != nil {
		return time.Time{}, err
	}

	ttlSeconds := int64(s.ttl.Seconds())
	pipe := s.redis.Pipeline()
	userKey := UserHoldsKey(userID, showtimeID)
	pipe.Expire(ctx, userKey, s.ttl)
	for _, seatID := range holds {
		pipe.Expire(ctx, SeatKey(showtimeID, seatID), s.ttl)
	}
	if _, err := pipe.Exec(ctx); err != nil {
		return time.Time{}, fmt.Errorf("refresh hold ttl: %w", err)
	}

	return s.now().Add(time.Duration(ttlSeconds) * time.Second), nil
}

func uniqueSorted(ids []string) []string {
	if len(ids) == 0 {
		return []string{}
	}
	seen := make(map[string]struct{}, len(ids))
	out := make([]string, 0, len(ids))
	for _, id := range ids {
		if id == "" {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		out = append(out, id)
	}
	slices.Sort(out)
	return out
}

func toSet(ids []string) map[string]struct{} {
	set := make(map[string]struct{}, len(ids))
	for _, id := range ids {
		set[id] = struct{}{}
	}
	return set
}
