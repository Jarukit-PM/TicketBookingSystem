package booking_test

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/booking"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/catalog"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/hold"
)

type stubShowtimes struct{ showtime *catalog.Showtime }

func (s stubShowtimes) InsertShowtime(context.Context, *catalog.Showtime) error { return nil }
func (s stubShowtimes) FindShowtimeByID(_ context.Context, id primitive.ObjectID) (*catalog.Showtime, error) {
	if s.showtime != nil && s.showtime.ID == id {
		return s.showtime, nil
	}
	return nil, nil
}
func (s stubShowtimes) ListShowtimesByScreen(context.Context, primitive.ObjectID, time.Time) ([]catalog.Showtime, error) {
	return nil, nil
}
func (s stubShowtimes) ListShowtimesByMovie(context.Context, primitive.ObjectID) ([]catalog.Showtime, error) {
	return nil, nil
}
func (s stubShowtimes) ListShowtimesByScreens(context.Context, []primitive.ObjectID) ([]catalog.Showtime, error) {
	return nil, nil
}
func (s stubShowtimes) ListShowtimesByCinemaMovie(context.Context, []primitive.ObjectID, primitive.ObjectID) ([]catalog.Showtime, error) {
	return nil, nil
}
func (stubShowtimes) ListAdminShowtimes(context.Context, catalog.AdminShowtimeFilter) ([]catalog.Showtime, error) {
	return nil, nil
}
func (stubShowtimes) UpdateShowtime(context.Context, *catalog.Showtime) error { return nil }
func (stubShowtimes) DeleteShowtime(context.Context, primitive.ObjectID) error { return nil }

type stubScreens struct{ screen *catalog.Screen }

func (s stubScreens) InsertScreen(context.Context, *catalog.Screen) error { return nil }
func (s stubScreens) FindScreenByID(_ context.Context, id primitive.ObjectID) (*catalog.Screen, error) {
	if s.screen != nil && s.screen.ID == id {
		return s.screen, nil
	}
	return nil, nil
}
func (s stubScreens) ListScreensByCinema(context.Context, primitive.ObjectID) ([]catalog.Screen, error) {
	return nil, nil
}
func (stubScreens) ListScreens(context.Context, *primitive.ObjectID) ([]catalog.Screen, error) {
	return nil, nil
}
func (stubScreens) UpdateScreen(context.Context, *catalog.Screen) error { return nil }
func (stubScreens) DeleteScreen(context.Context, primitive.ObjectID) error { return nil }

type stubCinemas struct{ cinema *catalog.Cinema }

func (s stubCinemas) InsertCinema(context.Context, *catalog.Cinema) error { return nil }
func (s stubCinemas) FindCinemaByID(_ context.Context, id primitive.ObjectID) (*catalog.Cinema, error) {
	if s.cinema != nil && s.cinema.ID == id {
		return s.cinema, nil
	}
	return nil, nil
}
func (s stubCinemas) ListCinemas(context.Context) ([]catalog.Cinema, error) { return nil, nil }
func (stubCinemas) UpdateCinema(context.Context, *catalog.Cinema) error      { return nil }
func (stubCinemas) DeleteCinema(context.Context, primitive.ObjectID) error { return nil }

type stubMovies struct{}

func (stubMovies) InsertMovie(context.Context, *catalog.Movie) error { return nil }
func (stubMovies) FindMovieByID(context.Context, primitive.ObjectID) (*catalog.Movie, error) {
	return nil, nil
}
func (stubMovies) ListMoviesByStatus(context.Context, string) ([]catalog.Movie, error) { return nil, nil }
func (stubMovies) ListComingSoonMovies(context.Context) ([]catalog.Movie, error)      { return nil, nil }
func (stubMovies) ListNonArchivedMovies(context.Context) ([]catalog.Movie, error)     { return nil, nil }
func (stubMovies) ListMovies(context.Context) ([]catalog.Movie, error)                { return nil, nil }
func (stubMovies) UpdateMovie(context.Context, *catalog.Movie) error                  { return nil }
func (stubMovies) DeleteMovie(context.Context, primitive.ObjectID) error              { return nil }

type memBookings struct {
	mu       sync.Mutex
	bookings []booking.Booking
	refs     map[string]struct{}
}

func newMemBookings() *memBookings {
	return &memBookings{refs: make(map[string]struct{})}
}

func (m *memBookings) Insert(_ context.Context, b *booking.Booking) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, exists := m.refs[b.BookingRef]; exists {
		return fmt.Errorf("duplicate key")
	}
	if b.ID.IsZero() {
		b.ID = primitive.NewObjectID()
	}
	m.refs[b.BookingRef] = struct{}{}
	m.bookings = append(m.bookings, *b)
	return nil
}

func (m *memBookings) FindByID(context.Context, primitive.ObjectID) (*booking.Booking, error) {
	return nil, nil
}

func (m *memBookings) FindByBookingRef(context.Context, string) (*booking.Booking, error) {
	return nil, nil
}

func (m *memBookings) UpdateTicketToken(_ context.Context, id primitive.ObjectID, token string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for i := range m.bookings {
		if m.bookings[i].ID == id {
			m.bookings[i].TicketToken = token
			return nil
		}
	}
	return nil
}

func (m *memBookings) ListByUser(ctx context.Context, userID primitive.ObjectID) ([]booking.Booking, error) {
	return m.ListConfirmedByUser(ctx, userID)
}

func (m *memBookings) ListConfirmedByUser(_ context.Context, userID primitive.ObjectID) ([]booking.Booking, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make([]booking.Booking, 0)
	for _, b := range m.bookings {
		if b.UserID == userID && b.Status == booking.StatusConfirmed {
			out = append(out, b)
		}
	}
	return out, nil
}

func (m *memBookings) CountConfirmedBetween(_ context.Context, from, to time.Time) (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	count := 0
	for _, b := range m.bookings {
		if b.Status == booking.StatusConfirmed && !b.ConfirmedAt.Before(from) && b.ConfirmedAt.Before(to) {
			count++
		}
	}
	return count, nil
}

func (m *memBookings) ListRecentConfirmed(_ context.Context, limit int) ([]booking.Booking, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make([]booking.Booking, 0, limit)
	for i := len(m.bookings) - 1; i >= 0 && len(out) < limit; i-- {
		if m.bookings[i].Status == booking.StatusConfirmed {
			out = append(out, m.bookings[i])
		}
	}
	return out, nil
}

func (m *memBookings) ListConfirmedByShowtime(_ context.Context, showtimeID primitive.ObjectID) ([]booking.Booking, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make([]booking.Booking, 0)
	for _, b := range m.bookings {
		if b.ShowtimeID == showtimeID && b.Status == booking.StatusConfirmed {
			out = append(out, b)
		}
	}
	return out, nil
}

func (m *memBookings) CountConfirmed(_ context.Context) (int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var count int64
	for _, b := range m.bookings {
		if b.Status == booking.StatusConfirmed {
			count++
		}
	}
	return count, nil
}

func (m *memBookings) CountConfirmedFiltered(_ context.Context, filter booking.ConfirmedFilter) (int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var count int64
	for _, b := range m.bookings {
		if filter.Matches(b) {
			count++
		}
	}
	return count, nil
}

func (m *memBookings) ListConfirmedFiltered(
	_ context.Context,
	filter booking.ConfirmedFilter,
	skip, limit int,
) ([]booking.Booking, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	matched := make([]booking.Booking, 0)
	for _, b := range m.bookings {
		if filter.Matches(b) {
			matched = append(matched, b)
		}
	}
	if skip >= len(matched) {
		return []booking.Booking{}, nil
	}
	end := skip + limit
	if end > len(matched) {
		end = len(matched)
	}
	return matched[skip:end], nil
}

func (m *memBookings) ListConfirmedPage(_ context.Context, skip, limit int) ([]booking.Booking, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	confirmed := make([]booking.Booking, 0)
	for _, b := range m.bookings {
		if b.Status == booking.StatusConfirmed {
			confirmed = append(confirmed, b)
		}
	}
	if skip >= len(confirmed) {
		return []booking.Booking{}, nil
	}
	end := skip + limit
	if end > len(confirmed) {
		end = len(confirmed)
	}
	return confirmed[skip:end], nil
}

type confirmEnv struct {
	mr       *miniredis.Miniredis
	rdb      *redis.Client
	holdSvc  *hold.Service
	bookSvc  *booking.Service
	bookings *memBookings
	showID   string
	userA    string
	userB    string
	userOID  primitive.ObjectID
}

func newConfirmEnv(t *testing.T, now time.Time, sold []string) *confirmEnv {
	t.Helper()

	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("miniredis: %v", err)
	}

	rdb := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	cinemaID := primitive.NewObjectID()
	screenID := primitive.NewObjectID()
	showtimeID := primitive.NewObjectID()
	userOID := primitive.NewObjectID()

	layoutSeats := []catalog.LayoutSeat{
		{SeatID: "A-1", Row: 1, Col: 1, Type: catalog.SeatTypeStandard},
		{SeatID: "A-2", Row: 1, Col: 2, Type: catalog.SeatTypeVIP},
		{SeatID: "A-3", Row: 1, Col: 3, Type: catalog.SeatTypeStandard},
	}

	showtime := &catalog.Showtime{
		ID:       showtimeID,
		ScreenID: screenID,
		StartsAt: now.Add(2 * time.Hour),
		Status:   catalog.ShowtimeStatusOpen,
		PriceTiers: catalog.PriceTiers{
			Standard: 1000,
			VIP:      1500,
		},
	}
	screen := &catalog.Screen{
		ID:       screenID,
		CinemaID: cinemaID,
		Name:     "Hall 1",
		Layout:   catalog.ScreenLayout{Seats: layoutSeats},
	}
	cinema := &catalog.Cinema{ID: cinemaID, Timezone: "UTC"}

	bookings := newMemBookings()
	if len(sold) > 0 {
		bookings.bookings = []booking.Booking{{
			ShowtimeID: showtimeID,
			Seats:      sold,
			Status:     booking.StatusConfirmed,
		}}
	}

	holdSvc := hold.NewService(
		stubShowtimes{showtime},
		stubScreens{screen},
		stubCinemas{cinema},
		bookings,
		rdb,
		hold.WithClock(func() time.Time { return now }),
	)
	bookSvc := booking.NewService(
		stubShowtimes{showtime},
		stubScreens{screen},
		stubCinemas{cinema},
		stubMovies{},
		bookings,
		holdSvc,
		rdb,
		booking.NewIdempotencyStore(rdb, time.Hour),
		booking.WithClock(func() time.Time { return now }),
		booking.WithTicketConfig("test-ticket-secret", "http://localhost"),
	)

	return &confirmEnv{
		mr:       mr,
		rdb:      rdb,
		holdSvc:  holdSvc,
		bookSvc:  bookSvc,
		bookings: bookings,
		showID:   showtimeID.Hex(),
		userA:    userOID.Hex(),
		userB:    primitive.NewObjectID().Hex(),
		userOID:  userOID,
	}
}

func (e *confirmEnv) close(t *testing.T) {
	t.Helper()
	e.rdb.Close()
	e.mr.Close()
}

func (e *confirmEnv) holdSeats(t *testing.T, ctx context.Context, userID string, seatIDs ...string) {
	t.Helper()
	if _, err := e.holdSvc.AddSeats(ctx, userID, e.showID, seatIDs); err != nil {
		t.Fatalf("hold seats: %v", err)
	}
}

func (e *confirmEnv) forceHold(t *testing.T, ctx context.Context, userID string, seatIDs ...string) {
	t.Helper()
	for _, seatID := range seatIDs {
		payload, err := json.Marshal(hold.Record{UserID: userID, HeldAt: time.Now().UTC()})
		if err != nil {
			t.Fatalf("marshal hold: %v", err)
		}
		if err := e.rdb.Set(ctx, hold.SeatKey(e.showID, seatID), payload, hold.DefaultTTL).Err(); err != nil {
			t.Fatalf("set hold: %v", err)
		}
		if err := e.rdb.SAdd(ctx, hold.UserHoldsKey(userID, e.showID), seatID).Err(); err != nil {
			t.Fatalf("track hold: %v", err)
		}
	}
}

func TestConfirm_IdempotencyHitReturnsSameBooking(t *testing.T) {
	env := newConfirmEnv(t, time.Date(2026, 6, 11, 18, 0, 0, 0, time.UTC), nil)
	defer env.close(t)

	ctx := context.Background()
	env.holdSeats(t, ctx, env.userA, "A-1", "A-2")

	key := "idem-key-1"
	first, err := env.bookSvc.Confirm(ctx, env.userA, env.showID, key, booking.LocaleEN)
	if err != nil {
		t.Fatalf("first confirm: %v", err)
	}

	second, err := env.bookSvc.Confirm(ctx, env.userA, env.showID, key, booking.LocaleEN)
	if err != nil {
		t.Fatalf("second confirm: %v", err)
	}

	if first.ID != second.ID {
		t.Fatalf("expected same booking id, got %s vs %s", first.ID.Hex(), second.ID.Hex())
	}
	if len(env.bookings.bookings) != 1 {
		t.Fatalf("expected 1 booking in store, got %d", len(env.bookings.bookings))
	}
}

func TestConfirm_NoHoldsReturnsConflict(t *testing.T) {
	env := newConfirmEnv(t, time.Date(2026, 6, 11, 18, 0, 0, 0, time.UTC), nil)
	defer env.close(t)

	_, err := env.bookSvc.Confirm(context.Background(), env.userA, env.showID, "fresh-key", booking.LocaleEN)
	if err == nil {
		t.Fatal("expected error")
	}
	if err != booking.ErrNoActiveHolds {
		t.Fatalf("expected ErrNoActiveHolds, got %v", err)
	}
}

func TestConfirm_ExpiredHoldsOnRetryReturnsConflict(t *testing.T) {
	now := time.Date(2026, 6, 11, 18, 0, 0, 0, time.UTC)
	env := newConfirmEnv(t, now, nil)
	defer env.close(t)

	ctx := context.Background()
	env.holdSeats(t, ctx, env.userA, "A-1")

	key := "retry-key"
	env.mr.FastForward(6 * time.Minute)

	_, err := env.bookSvc.Confirm(ctx, env.userA, env.showID, key, booking.LocaleEN)
	if err == nil {
		t.Fatal("expected error on confirm without active holds")
	}
	if err != booking.ErrNoActiveHolds {
		t.Fatalf("expected ErrNoActiveHolds, got %v", err)
	}
}

func TestConfirm_IdempotencyReturnsCachedAfterHoldsExpire(t *testing.T) {
	now := time.Date(2026, 6, 11, 18, 0, 0, 0, time.UTC)
	env := newConfirmEnv(t, now, nil)
	defer env.close(t)

	ctx := context.Background()
	env.holdSeats(t, ctx, env.userA, "A-1")

	key := "cached-key"
	first, err := env.bookSvc.Confirm(ctx, env.userA, env.showID, key, booking.LocaleEN)
	if err != nil {
		t.Fatalf("confirm: %v", err)
	}

	env.mr.FastForward(6 * time.Minute)

	second, err := env.bookSvc.Confirm(ctx, env.userA, env.showID, key, booking.LocaleEN)
	if err != nil {
		t.Fatalf("idempotent retry: %v", err)
	}
	if first.ID != second.ID {
		t.Fatalf("expected cached booking %s, got %s", first.ID.Hex(), second.ID.Hex())
	}
}

func TestConfirm_TotalPriceAndBookingRef(t *testing.T) {
	env := newConfirmEnv(t, time.Date(2026, 6, 11, 18, 0, 0, 0, time.UTC), nil)
	defer env.close(t)

	ctx := context.Background()
	env.holdSeats(t, ctx, env.userA, "A-1", "A-2")

	got, err := env.bookSvc.Confirm(ctx, env.userA, env.showID, "price-key", booking.LocaleEN)
	if err != nil {
		t.Fatalf("confirm: %v", err)
	}

	wantTotal := int64(2500)
	if got.Total != wantTotal {
		t.Fatalf("total = %d, want %d", got.Total, wantTotal)
	}
	if !booking.ValidateBookingRefFormat(got.BookingRef) {
		t.Fatalf("invalid booking ref: %q", got.BookingRef)
	}
	if got.TicketToken == "" || got.TicketToken == got.BookingRef {
		t.Fatalf("expected opaque ticket token distinct from booking ref")
	}
	if !booking.ValidateTicketToken(got.BookingRef, got.TicketToken, got, "test-ticket-secret") {
		t.Fatalf("expected signed ticket token to validate")
	}
}

func TestConfirm_MultipleBookingsSameUserShowtime(t *testing.T) {
	env := newConfirmEnv(t, time.Date(2026, 6, 11, 18, 0, 0, 0, time.UTC), nil)
	defer env.close(t)

	ctx := context.Background()
	env.holdSeats(t, ctx, env.userA, "A-1")
	if _, err := env.bookSvc.Confirm(ctx, env.userA, env.showID, "booking-1", booking.LocaleEN); err != nil {
		t.Fatalf("first confirm: %v", err)
	}

	env.holdSeats(t, ctx, env.userA, "A-2")
	second, err := env.bookSvc.Confirm(ctx, env.userA, env.showID, "booking-2", booking.LocaleEN)
	if err != nil {
		t.Fatalf("second confirm: %v", err)
	}

	if len(env.bookings.bookings) != 2 {
		t.Fatalf("expected 2 bookings, got %d", len(env.bookings.bookings))
	}
	if len(second.Seats) != 1 || second.Seats[0] != "A-2" {
		t.Fatalf("unexpected second booking seats: %v", second.Seats)
	}
}

func TestConfirm_ConcurrentDuplicateConfirmOnlyOneBooking(t *testing.T) {
	env := newConfirmEnv(t, time.Date(2026, 6, 11, 18, 0, 0, 0, time.UTC), nil)
	defer env.close(t)

	ctx := context.Background()
	env.holdSeats(t, ctx, env.userA, "A-1")

	var wg sync.WaitGroup
	errs := make([]error, 2)
	wg.Add(2)
	for i := range 2 {
		go func(idx int) {
			defer wg.Done()
			_, errs[idx] = env.bookSvc.Confirm(ctx, env.userA, env.showID, fmt.Sprintf("race-key-%d", idx), booking.LocaleEN)
		}(i)
	}
	wg.Wait()

	successes := 0
	for _, err := range errs {
		if err == nil {
			successes++
			continue
		}
		if err != booking.ErrNoActiveHolds && err != booking.ErrSeatConflict {
			t.Fatalf("unexpected error: %v", err)
		}
	}
	if successes != 1 {
		t.Fatalf("expected exactly 1 success, got %d", successes)
	}
	if len(env.bookings.bookings) != 1 {
		t.Fatalf("expected 1 persisted booking, got %d", len(env.bookings.bookings))
	}
}

func TestConfirm_SecondUserWithoutHoldGetsConflict(t *testing.T) {
	env := newConfirmEnv(t, time.Date(2026, 6, 11, 18, 0, 0, 0, time.UTC), nil)
	defer env.close(t)

	ctx := context.Background()
	env.holdSeats(t, ctx, env.userA, "A-1")

	_, err := env.bookSvc.Confirm(ctx, env.userB, env.showID, "user-b-key", booking.LocaleEN)
	if err == nil {
		t.Fatal("expected error")
	}
	if err != booking.ErrNoActiveHolds {
		t.Fatalf("expected ErrNoActiveHolds, got %v", err)
	}
}

func TestConfirm_ClearsHoldsOnSuccess(t *testing.T) {
	env := newConfirmEnv(t, time.Date(2026, 6, 11, 18, 0, 0, 0, time.UTC), nil)
	defer env.close(t)

	ctx := context.Background()
	env.holdSeats(t, ctx, env.userA, "A-1", "A-2")

	if _, err := env.bookSvc.Confirm(ctx, env.userA, env.showID, "clear-key", booking.LocaleEN); err != nil {
		t.Fatalf("confirm: %v", err)
	}

	holds, err := env.holdSvc.UserHolds(ctx, env.userA, env.showID)
	if err != nil {
		t.Fatalf("list holds: %v", err)
	}
	if len(holds) != 0 {
		t.Fatalf("expected holds cleared, got %v", holds)
	}

	raw, err := env.rdb.Get(ctx, hold.SeatKey(env.showID, "A-1")).Result()
	if err != redis.Nil {
		t.Fatalf("expected hold key gone, got %q err=%v", raw, err)
	}
}

func TestConfirm_RejectsAlreadySoldSeat(t *testing.T) {
	env := newConfirmEnv(t, time.Date(2026, 6, 11, 18, 0, 0, 0, time.UTC), []string{"A-1"})
	defer env.close(t)

	ctx := context.Background()
	env.forceHold(t, ctx, env.userA, "A-1")

	_, err := env.bookSvc.Confirm(ctx, env.userA, env.showID, "sold-key", booking.LocaleEN)
	if err == nil {
		t.Fatal("expected error")
	}
	if err != booking.ErrSeatConflict {
		t.Fatalf("expected ErrSeatConflict, got %v", err)
	}
}

func TestConfirm_HoldPayloadStillOwned(t *testing.T) {
	env := newConfirmEnv(t, time.Date(2026, 6, 11, 18, 0, 0, 0, time.UTC), nil)
	defer env.close(t)

	ctx := context.Background()
	env.holdSeats(t, ctx, env.userA, "A-1")

	raw, err := env.rdb.Get(ctx, hold.SeatKey(env.showID, "A-1")).Result()
	if err != nil {
		t.Fatalf("get hold: %v", err)
	}
	var rec hold.Record
	if err := json.Unmarshal([]byte(raw), &rec); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if rec.UserID != env.userA {
		t.Fatalf("hold owner = %q, want %q", rec.UserID, env.userA)
	}
}

func TestConfirm_PersistsLocale(t *testing.T) {
	env := newConfirmEnv(t, time.Date(2026, 6, 11, 18, 0, 0, 0, time.UTC), nil)
	defer env.close(t)

	ctx := context.Background()
	env.holdSeats(t, ctx, env.userA, "A-1")

	got, err := env.bookSvc.Confirm(ctx, env.userA, env.showID, "locale-key", booking.LocaleTH)
	if err != nil {
		t.Fatalf("confirm: %v", err)
	}
	if got.Locale != booking.LocaleTH {
		t.Fatalf("locale = %q, want %q", got.Locale, booking.LocaleTH)
	}
}
