package hold_test

import (
	"context"
	"encoding/json"
	"errors"
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

type stubBookings struct{ bookings []booking.Booking }

func (s stubBookings) Insert(context.Context, *booking.Booking) error { return nil }
func (s stubBookings) FindByID(context.Context, primitive.ObjectID) (*booking.Booking, error) {
	return nil, nil
}
func (s stubBookings) FindByBookingRef(context.Context, string) (*booking.Booking, error) {
	return nil, nil
}
func (s stubBookings) ListByUser(ctx context.Context, userID primitive.ObjectID) ([]booking.Booking, error) {
	return s.ListConfirmedByUser(ctx, userID)
}
func (s stubBookings) ListConfirmedByUser(context.Context, primitive.ObjectID) ([]booking.Booking, error) {
	return nil, nil
}
func (stubBookings) CountConfirmedBetween(context.Context, time.Time, time.Time) (int, error) {
	return 0, nil
}
func (stubBookings) ListRecentConfirmed(context.Context, int) ([]booking.Booking, error) {
	return nil, nil
}
func (s stubBookings) ListConfirmedByShowtime(_ context.Context, showtimeID primitive.ObjectID) ([]booking.Booking, error) {
	out := make([]booking.Booking, 0)
	for _, b := range s.bookings {
		if b.ShowtimeID == showtimeID {
			out = append(out, b)
		}
	}
	return out, nil
}

type testEnv struct {
	mr     *miniredis.Miniredis
	rdb    *redis.Client
	svc    *hold.Service
	showID string
	userA  string
	userB  string
}

func newTestEnv(t *testing.T, startsAt time.Time, sold []string, now time.Time) *testEnv {
	t.Helper()

	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("miniredis: %v", err)
	}

	rdb := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	cinemaID := primitive.NewObjectID()
	screenID := primitive.NewObjectID()
	showtimeID := primitive.NewObjectID()

	layoutSeats := make([]catalog.LayoutSeat, 0, 13)
	for i := 1; i <= 12; i++ {
		layoutSeats = append(layoutSeats, catalog.LayoutSeat{
			SeatID: fmt.Sprintf("A-%d", i), Row: 1, Col: i, Type: catalog.SeatTypeStandard,
		})
	}
	layoutSeats = append(layoutSeats, catalog.LayoutSeat{SeatID: "B-1", Row: 2, Col: 1, Type: catalog.SeatTypeBlocked})

	showtime := &catalog.Showtime{ID: showtimeID, ScreenID: screenID, StartsAt: startsAt, Status: catalog.ShowtimeStatusOpen}
	screen := &catalog.Screen{ID: screenID, CinemaID: cinemaID, Name: "Hall 1", Layout: catalog.ScreenLayout{Seats: layoutSeats}}
	cinema := &catalog.Cinema{ID: cinemaID, Timezone: "UTC"}

	bookings := stubBookings{}
	if len(sold) > 0 {
		bookings.bookings = []booking.Booking{{ShowtimeID: showtimeID, Seats: sold, Status: booking.StatusConfirmed}}
	}

	svc := hold.NewService(stubShowtimes{showtime}, stubScreens{screen}, stubCinemas{cinema}, bookings, rdb, hold.WithClock(func() time.Time { return now }))

	return &testEnv{mr: mr, rdb: rdb, svc: svc, showID: showtimeID.Hex(), userA: "user-a", userB: "user-b"}
}

func (e *testEnv) close(t *testing.T) {
	t.Helper()
	e.rdb.Close()
	e.mr.Close()
}

func TestAddSeats_holdRules(t *testing.T) {
	now := time.Date(2026, 6, 11, 18, 0, 0, 0, time.UTC)
	future := now.Add(2 * time.Hour)
	past := now.Add(-1 * time.Hour)

	tests := []struct {
		name    string
		starts  time.Time
		sold    []string
		setup   func(e *testEnv, ctx context.Context)
		userID  string
		add     []string
		wantErr error
		assert  func(t *testing.T, e *testEnv, ctx context.Context, res hold.Result)
	}{
		{
			name: "adds seat and sets redis keys with ttl", starts: future, userID: "user-a", add: []string{"A-1"},
			assert: func(t *testing.T, e *testEnv, _ context.Context, res hold.Result) {
				if len(res.Holds) != 1 || res.Holds[0] != "A-1" || res.ExpiresAt == nil {
					t.Fatalf("unexpected result: %+v", res)
				}
				if e.mr.TTL(hold.SeatKey(e.showID, "A-1")) != 5*time.Minute {
					t.Fatal("seat ttl not 5m")
				}
			},
		},
		{
			name: "refreshes ttl on add not on remove", starts: future, userID: "user-a", add: []string{"A-2"},
			setup: func(e *testEnv, ctx context.Context) {
				if _, err := e.svc.AddSeats(ctx, e.userA, e.showID, []string{"A-1"}); err != nil {
					t.Fatalf("setup: %v", err)
				}
				e.mr.FastForward(4 * time.Minute)
			},
			assert: func(t *testing.T, e *testEnv, ctx context.Context, _ hold.Result) {
				if e.mr.TTL(hold.SeatKey(e.showID, "A-1")) < 4*time.Minute+30*time.Second {
					t.Fatal("A-1 ttl not refreshed on add")
				}
				before := e.mr.TTL(hold.SeatKey(e.showID, "A-1"))
				if _, err := e.svc.RemoveSeats(ctx, e.userA, e.showID, []string{"A-2"}); err != nil {
					t.Fatalf("remove: %v", err)
				}
				e.mr.FastForward(30 * time.Second)
				if e.mr.TTL(hold.SeatKey(e.showID, "A-1")) >= before {
					t.Fatal("ttl should not refresh on remove")
				}
			},
		},
		{
			name: "rejects more than ten seats", starts: future, userID: "user-a",
			add: []string{"A-1", "A-2", "A-3", "A-4", "A-5", "A-6", "A-7", "A-8", "A-9", "A-10", "A-11"},
			wantErr: hold.ErrSeatLimitExceeded,
		},
		{name: "rejects sold seat", starts: future, sold: []string{"A-1"}, userID: "user-a", add: []string{"A-1"}, wantErr: hold.ErrSeatSold},
		{name: "rejects blocked seat", starts: future, userID: "user-a", add: []string{"B-1"}, wantErr: hold.ErrSeatBlocked},
		{name: "rejects after showtime started", starts: past, userID: "user-a", add: []string{"A-1"}, wantErr: hold.ErrShowtimeStarted},
		{
			name: "rejects seat held by another user", starts: future, userID: "user-b", add: []string{"A-1"},
			setup: func(e *testEnv, ctx context.Context) {
				if _, err := e.svc.AddSeats(ctx, e.userA, e.showID, []string{"A-1"}); err != nil {
					t.Fatalf("setup: %v", err)
				}
			},
			wantErr: hold.ErrSeatHeldByOther,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := newTestEnv(t, tt.starts, tt.sold, now)
			defer env.close(t)
			ctx := context.Background()
			if tt.setup != nil {
				tt.setup(env, ctx)
			}
			userID := tt.userID
			if userID == "" {
				userID = env.userA
			}
			res, err := env.svc.AddSeats(ctx, userID, env.showID, tt.add)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("err = %v, want %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("AddSeats: %v", err)
			}
			if tt.assert != nil {
				tt.assert(t, env, ctx, res)
			}
		})
	}
}

func TestRemoveSeats_releaseAll(t *testing.T) {
	now := time.Date(2026, 6, 11, 18, 0, 0, 0, time.UTC)
	env := newTestEnv(t, now.Add(2*time.Hour), nil, now)
	defer env.close(t)
	ctx := context.Background()

	if _, err := env.svc.AddSeats(ctx, env.userA, env.showID, []string{"A-1", "A-2"}); err != nil {
		t.Fatalf("add: %v", err)
	}
	res, err := env.svc.RemoveSeats(ctx, env.userA, env.showID, nil)
	if err != nil || len(res.Holds) != 0 {
		t.Fatalf("remove all: err=%v holds=%v", err, res.Holds)
	}
}

func TestAddSeats_setNXRaceOnlyOneWinner(t *testing.T) {
	now := time.Date(2026, 6, 11, 18, 0, 0, 0, time.UTC)
	env := newTestEnv(t, now.Add(2*time.Hour), nil, now)
	defer env.close(t)
	ctx := context.Background()

	var wg sync.WaitGroup
	errs := make(chan error, 2)
	wg.Add(2)
	go func() { defer wg.Done(); _, err := env.svc.AddSeats(ctx, env.userA, env.showID, []string{"A-1"}); errs <- err }()
	go func() { defer wg.Done(); _, err := env.svc.AddSeats(ctx, env.userB, env.showID, []string{"A-1"}); errs <- err }()
	wg.Wait()
	close(errs)

	winners, losers := 0, 0
	for err := range errs {
		switch {
		case err == nil:
			winners++
		case errors.Is(err, hold.ErrSeatHeldByOther):
			losers++
		default:
			t.Fatalf("unexpected err: %v", err)
		}
	}
	if winners != 1 || losers != 1 {
		t.Fatalf("winners=%d losers=%d", winners, losers)
	}

	raw, _ := env.rdb.Get(ctx, hold.SeatKey(env.showID, "A-1")).Result()
	var rec hold.Record
	_ = json.Unmarshal([]byte(raw), &rec)
	if rec.UserID != env.userA && rec.UserID != env.userB {
		t.Fatalf("unexpected owner %s", rec.UserID)
	}
}
