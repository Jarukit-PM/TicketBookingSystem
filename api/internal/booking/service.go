package booking

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"time"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/catalog"
)

const maxSeatsPerBooking = 10

// HoldStore lists and clears Redis seat holds for confirm.
type HoldStore interface {
	UserHolds(ctx context.Context, userID, showtimeID string) ([]string, error)
	ClearUserHolds(ctx context.Context, userID, showtimeID string) error
}

// Service confirms bookings from active seat holds.
type Service struct {
	showtimes   catalog.ShowtimeRepository
	screens     catalog.ScreenRepository
	cinemas     catalog.CinemaRepository
	movies      catalog.MovieRepository
	bookings    Repository
	holds       HoldStore
	redis       *redis.Client
	idempotency *IdempotencyStore
	now         func() time.Time
}

// Option configures a booking Service.
type Option func(*Service)

// WithClock overrides the clock used for cutoff checks.
func WithClock(now func() time.Time) Option {
	return func(s *Service) { s.now = now }
}

// NewService returns a booking confirm service.
func NewService(
	showtimes catalog.ShowtimeRepository,
	screens catalog.ScreenRepository,
	cinemas catalog.CinemaRepository,
	movies catalog.MovieRepository,
	bookings Repository,
	holds HoldStore,
	rdb *redis.Client,
	idempotency *IdempotencyStore,
	opts ...Option,
) *Service {
	s := &Service{
		showtimes:   showtimes,
		screens:     screens,
		cinemas:     cinemas,
		movies:      movies,
		bookings:    bookings,
		holds:       holds,
		redis:       rdb,
		idempotency: idempotency,
		now:         time.Now,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

type confirmCtx struct {
	showtimeID string
	showtime   *catalog.Showtime
	screen     *catalog.Screen
	layout     map[string]catalog.LayoutSeat
	sold       map[string]struct{}
}

// Confirm books all held seats for one showtime as a single CONFIRMED booking.
func (s *Service) Confirm(ctx context.Context, userID, showtimeID, idempotencyKey string) (*Booking, error) {
	if idempotencyKey == "" {
		return nil, ErrIdempotencyRequired
	}

	if cached, err := s.idempotency.Get(ctx, idempotencyKey); err != nil {
		return nil, err
	} else if cached != nil {
		return cached, nil
	}

	seatIDs, err := s.holds.UserHolds(ctx, userID, showtimeID)
	if err != nil {
		return nil, err
	}
	seatIDs = uniqueSorted(seatIDs)
	if len(seatIDs) == 0 {
		return nil, ErrNoActiveHolds
	}
	if len(seatIDs) > maxSeatsPerBooking {
		return nil, ErrSeatLimitExceeded
	}

	stx, err := s.loadConfirmCtx(ctx, showtimeID)
	if err != nil {
		return nil, err
	}
	if err := s.ensureShowtimeOpen(ctx, stx); err != nil {
		return nil, err
	}

	for _, seatID := range seatIDs {
		if err := s.validateSeatForConfirm(ctx, stx, userID, showtimeID, seatID); err != nil {
			return nil, err
		}
	}

	releaseLocks, err := acquireConfirmLocks(ctx, s.redis, showtimeID, seatIDs)
	if err != nil {
		return nil, err
	}
	defer releaseLocks(ctx)

	stx, err = s.loadConfirmCtx(ctx, showtimeID)
	if err != nil {
		return nil, err
	}
	for _, seatID := range seatIDs {
		if _, sold := stx.sold[seatID]; sold {
			return nil, ErrSeatConflict
		}
		if err := s.validateSeatForConfirm(ctx, stx, userID, showtimeID, seatID); err != nil {
			return nil, err
		}
	}

	total, err := catalog.TotalForSeats(stx.showtime.PriceTiers, stx.screen.Layout.Seats, seatIDs)
	if err != nil {
		return nil, fmt.Errorf("price seats: %w", err)
	}

	userOID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user id: %w", err)
	}
	showtimeOID, err := primitive.ObjectIDFromHex(showtimeID)
	if err != nil {
		return nil, fmt.Errorf("invalid showtime id: %w", err)
	}

	confirmedAt := s.now().UTC()
	booking, err := s.insertBookingWithRetry(ctx, userOID, showtimeOID, seatIDs, total, confirmedAt)
	if err != nil {
		return nil, err
	}

	if err := s.holds.ClearUserHolds(ctx, userID, showtimeID); err != nil {
		return nil, err
	}

	if err := s.idempotency.Set(ctx, idempotencyKey, booking); err != nil {
		return nil, err
	}

	return booking, nil
}

func (s *Service) insertBookingWithRetry(
	ctx context.Context,
	userID, showtimeID primitive.ObjectID,
	seatIDs []string,
	total int64,
	confirmedAt time.Time,
) (*Booking, error) {
	const maxAttempts = 5
	for attempt := 0; attempt < maxAttempts; attempt++ {
		ref, err := GenerateBookingRef()
		if err != nil {
			return nil, err
		}
		token, err := GenerateTicketToken()
		if err != nil {
			return nil, err
		}

		b := &Booking{
			UserID:      userID,
			ShowtimeID:  showtimeID,
			Seats:       seatIDs,
			Total:       total,
			BookingRef:  ref,
			TicketToken: token,
			Status:      StatusConfirmed,
			ConfirmedAt: confirmedAt,
		}
		if err := s.bookings.Insert(ctx, b); err != nil {
			if isDuplicateKey(err) {
				continue
			}
			return nil, err
		}
		return b, nil
	}
	return nil, fmt.Errorf("insert booking: exhausted bookingRef retries")
}

func (s *Service) loadConfirmCtx(ctx context.Context, showtimeID string) (*confirmCtx, error) {
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
		return nil, fmt.Errorf("screen not found")
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

	return &confirmCtx{
		showtimeID: showtimeID,
		showtime:   showtime,
		screen:     screen,
		layout:     layout,
		sold:       sold,
	}, nil
}

func (s *Service) ensureShowtimeOpen(ctx context.Context, stx *confirmCtx) error {
	cinema, err := s.cinemas.FindCinemaByID(ctx, stx.screen.CinemaID)
	if err != nil {
		return fmt.Errorf("find cinema: %w", err)
	}
	if cinema == nil {
		return fmt.Errorf("cinema not found")
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

func (s *Service) validateSeatForConfirm(
	ctx context.Context,
	stx *confirmCtx,
	userID, showtimeID, seatID string,
) error {
	seat, ok := stx.layout[seatID]
	if !ok {
		return ErrSeatConflict
	}
	if seat.Type == catalog.SeatTypeBlocked {
		return ErrSeatConflict
	}
	if _, sold := stx.sold[seatID]; sold {
		return ErrSeatConflict
	}

	owner, err := holdOwner(ctx, s.redis, showtimeID, seatID)
	if err != nil {
		return err
	}
	if owner != userID {
		return ErrSeatConflict
	}
	return nil
}

type holdRecord struct {
	UserID string `json:"userId"`
}

func seatHoldKey(showtimeID, seatID string) string {
	return fmt.Sprintf("hold:%s:%s", showtimeID, seatID)
}

func holdOwner(ctx context.Context, rdb *redis.Client, showtimeID, seatID string) (string, error) {
	raw, err := rdb.Get(ctx, seatHoldKey(showtimeID, seatID)).Result()
	if err == redis.Nil {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("get hold: %w", err)
	}

	var rec holdRecord
	if err := json.Unmarshal([]byte(raw), &rec); err != nil {
		return "", fmt.Errorf("decode hold: %w", err)
	}
	return rec.UserID, nil
}

func isDuplicateKey(err error) bool {
	return mongo.IsDuplicateKeyError(err)
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
