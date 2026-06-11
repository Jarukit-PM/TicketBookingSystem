package booking_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/booking"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/catalog"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/user"
)

var queryNow = time.Date(2026, 6, 11, 18, 0, 0, 0, time.UTC)

func newQueryEnv(t *testing.T, startsAt time.Time) (*booking.Service, primitive.ObjectID, primitive.ObjectID) {
	t.Helper()
	cinemaID := primitive.NewObjectID()
	screenID := primitive.NewObjectID()
	showtimeID := primitive.NewObjectID()
	movieID := primitive.NewObjectID()
	userID := primitive.NewObjectID()
	otherID := primitive.NewObjectID()
	movie := &catalog.Movie{ID: movieID, Title: "Query Movie", PosterURL: "https://example.com/poster.jpg"}
	showtime := &catalog.Showtime{ID: showtimeID, MovieID: movieID, ScreenID: screenID, StartsAt: startsAt, Status: catalog.ShowtimeStatusOpen}
	screen := &catalog.Screen{ID: screenID, CinemaID: cinemaID, Name: "Hall 1"}
	cinema := &catalog.Cinema{ID: cinemaID, Name: "Central Cinema", Timezone: "UTC"}
	bookings := newMemBookings()
	bookings.bookings = []booking.Booking{
		{ID: primitive.NewObjectID(), UserID: userID, ShowtimeID: showtimeID, Seats: []string{"A-1"}, Total: 1000, BookingRef: "TBS-QUERY1", Status: booking.StatusConfirmed, ConfirmedAt: queryNow.Add(-time.Hour)},
		{ID: primitive.NewObjectID(), UserID: userID, ShowtimeID: showtimeID, Seats: []string{"A-2"}, Total: 1000, BookingRef: "TBS-QUERY2", Status: booking.StatusConfirmed, ConfirmedAt: queryNow.Add(-30 * time.Minute)},
	}
	svc := booking.NewService(stubShowtimes{showtime}, stubScreens{screen}, stubCinemas{cinema}, stubMovies{movie}, bookings, nil, nil, nil, booking.WithClock(func() time.Time { return queryNow }))
	return svc, userID, otherID
}

func TestListMine_EnrichesBookings(t *testing.T) {
	svc, userID, _ := newQueryEnv(t, queryNow.Add(2*time.Hour))
	items, err := svc.ListMine(context.Background(), userID.Hex(), true)
	if err != nil || len(items) != 2 || items[0].Movie.Title != "Query Movie" {
		t.Fatalf("ListMine: err=%v items=%+v", err, items)
	}
}

func TestListMine_UpcomingFilter(t *testing.T) {
	svc, userID, _ := newQueryEnv(t, queryNow.Add(2*time.Hour))
	if upcoming, err := svc.ListMine(context.Background(), userID.Hex(), true); err != nil || len(upcoming) != 2 {
		t.Fatalf("upcoming: %v", err)
	}
	svcPast, userPast, _ := newQueryEnv(t, queryNow.Add(-2*time.Hour))
	if history, err := svcPast.ListMine(context.Background(), userPast.Hex(), false); err != nil || len(history) != 2 {
		t.Fatalf("history: %v", err)
	}
	if mixed, err := svcPast.ListMine(context.Background(), userPast.Hex(), true); err != nil || len(mixed) != 0 {
		t.Fatalf("past upcoming: %v len=%d", err, len(mixed))
	}
}

func TestGetDetail_OwnerAllowed(t *testing.T) {
	svc, userID, _ := newQueryEnv(t, queryNow.Add(2*time.Hour))
	items, _ := svc.ListMine(context.Background(), userID.Hex(), true)
	if _, err := svc.GetDetail(context.Background(), userID.Hex(), user.RoleCustomer, items[0].ID); err != nil {
		t.Fatalf("owner: %v", err)
	}
}

func TestGetDetail_ForbiddenForOtherUser(t *testing.T) {
	svc, userID, otherID := newQueryEnv(t, queryNow.Add(2*time.Hour))
	items, _ := svc.ListMine(context.Background(), userID.Hex(), true)
	if _, err := svc.GetDetail(context.Background(), otherID.Hex(), user.RoleCustomer, items[0].ID); !errors.Is(err, booking.ErrForbidden) {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}
}

func TestGetDetail_AdminAllowed(t *testing.T) {
	svc, userID, otherID := newQueryEnv(t, queryNow.Add(2*time.Hour))
	items, _ := svc.ListMine(context.Background(), userID.Hex(), true)
	if _, err := svc.GetDetail(context.Background(), otherID.Hex(), user.RoleAdmin, items[0].ID); err != nil {
		t.Fatalf("admin: %v", err)
	}
}

func TestListMine_MultipleBookingsSameShowtime(t *testing.T) {
	svc, userID, _ := newQueryEnv(t, queryNow.Add(2*time.Hour))
	items, err := svc.ListMine(context.Background(), userID.Hex(), true)
	if err != nil || len(items) != 2 || items[0].ShowtimeID != items[1].ShowtimeID {
		t.Fatalf("same showtime: %+v err=%v", items, err)
	}
}
