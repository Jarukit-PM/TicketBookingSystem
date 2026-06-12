package email_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/hibiken/asynq"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/audit"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/booking"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/catalog"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/email"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/tasks"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/user"
)

func TestHandleEmailSend(t *testing.T) {
	t.Parallel()

	bid := primitive.NewObjectID()
	uid := primitive.NewObjectID()
	sid := primitive.NewObjectID()
	scid := primitive.NewObjectID()
	mid := primitive.NewObjectID()
	cid := primitive.NewObjectID()
	b := &booking.Booking{
		ID: bid, UserID: uid, ShowtimeID: sid, Seats: []string{"A-1"}, Total: 1000,
		BookingRef: "TBS-T1", TicketToken: "tok", Status: booking.StatusConfirmed,
	}
	sender := &fakeSender{id: "sg-1"}
	logs := &fakeLogs{}
	svc := email.NewService(
		fakeBookings{b},
		fakeUsers{&user.User{ID: uid, Email: "a@b.com"}},
		email.CatalogReader{
			Showtimes: fakeST{&catalog.Showtime{ID: sid, ScreenID: scid, MovieID: mid, StartsAt: time.Now()}},
			Screens:   fakeSC{&catalog.Screen{ID: scid, CinemaID: cid, Name: "Hall"}},
			Movies:    fakeMV{&catalog.Movie{ID: mid, Title: "Film"}},
			Cinemas:   fakeCN{&catalog.Cinema{ID: cid, Name: "Cinema"}},
		},
		logs,
		sender,
		"http://localhost",
		"test-ticket-secret",
	)
	payload, _ := json.Marshal(tasks.EmailSendPayload{BookingID: bid.Hex()})
	if err := svc.HandleEmailSend(context.Background(), asynq.NewTask(tasks.TypeEmailSend, payload)); err != nil {
		t.Fatal(err)
	}
	if !sender.ok || logs.rows[0].Status != email.StatusSent {
		t.Fatal("expected sent email log")
	}
}

type fakeBookings struct{ b *booking.Booking }

func (f fakeBookings) Insert(context.Context, *booking.Booking) error { return nil }
func (f fakeBookings) FindByID(_ context.Context, id primitive.ObjectID) (*booking.Booking, error) {
	if f.b.ID == id {
		return f.b, nil
	}
	return nil, nil
}
func (f fakeBookings) FindByBookingRef(context.Context, string) (*booking.Booking, error) {
	return nil, nil
}
func (f fakeBookings) ListByUser(ctx context.Context, userID primitive.ObjectID) ([]booking.Booking, error) {
	return f.ListConfirmedByUser(ctx, userID)
}
func (f fakeBookings) ListConfirmedByUser(context.Context, primitive.ObjectID) ([]booking.Booking, error) {
	return nil, nil
}
func (fakeBookings) CountConfirmedBetween(context.Context, time.Time, time.Time) (int, error) {
	return 0, nil
}
func (fakeBookings) ListRecentConfirmed(context.Context, int) ([]booking.Booking, error) {
	return nil, nil
}
func (fakeBookings) CountConfirmed(context.Context) (int64, error) { return 0, nil }
func (fakeBookings) CountConfirmedFiltered(ctx context.Context, _ booking.ConfirmedFilter) (int64, error) {
	return 0, nil
}
func (fakeBookings) ListConfirmedPage(context.Context, int, int) ([]booking.Booking, error) {
	return nil, nil
}
func (fakeBookings) ListConfirmedFiltered(ctx context.Context, _ booking.ConfirmedFilter, _, _ int) ([]booking.Booking, error) {
	return nil, nil
}
func (f fakeBookings) ListConfirmedByShowtime(context.Context, primitive.ObjectID) ([]booking.Booking, error) {
	return nil, nil
}
func (f fakeBookings) UpdateTicketToken(_ context.Context, id primitive.ObjectID, token string) error {
	if f.b.ID == id {
		f.b.TicketToken = token
	}
	return nil
}

type fakeUsers struct{ u *user.User }

func (f fakeUsers) Insert(context.Context, *user.User) error { return nil }
func (f fakeUsers) FindByID(_ context.Context, id primitive.ObjectID) (*user.User, error) {
	if f.u.ID == id {
		return f.u, nil
	}
	return nil, nil
}
func (f fakeUsers) FindByEmail(context.Context, string) (*user.User, error) { return nil, nil }
func (f fakeUsers) FindByGoogleID(context.Context, string) (*user.User, error) {
	return nil, nil
}
func (f fakeUsers) Update(context.Context, *user.User) error { return nil }

type fakeSender struct {
	ok bool
	id string
}

func (f *fakeSender) Send(context.Context, email.Message) (string, error) {
	f.ok = true
	return f.id, nil
}

type fakeLogs struct{ rows []audit.EmailLog }

func (f *fakeLogs) InsertEmailLog(_ context.Context, l *audit.EmailLog) error {
	f.rows = append(f.rows, *l)
	return nil
}
func (f *fakeLogs) ListByBooking(context.Context, primitive.ObjectID) ([]audit.EmailLog, error) {
	return f.rows, nil
}
func (f *fakeLogs) ListEmailLogs(_ context.Context, _ audit.LogPage, _ *primitive.ObjectID) ([]audit.EmailLog, error) {
	return f.rows, nil
}

type fakeST struct{ s *catalog.Showtime }

func (f fakeST) InsertShowtime(context.Context, *catalog.Showtime) error { return nil }
func (f fakeST) FindShowtimeByID(context.Context, primitive.ObjectID) (*catalog.Showtime, error) {
	return f.s, nil
}
func (f fakeST) ListShowtimesByScreen(context.Context, primitive.ObjectID, time.Time) ([]catalog.Showtime, error) {
	return nil, nil
}
func (f fakeST) ListShowtimesByMovie(context.Context, primitive.ObjectID) ([]catalog.Showtime, error) {
	return nil, nil
}
func (f fakeST) ListShowtimesByScreens(context.Context, []primitive.ObjectID) ([]catalog.Showtime, error) {
	return nil, nil
}
func (f fakeST) ListShowtimesByCinemaMovie(context.Context, []primitive.ObjectID, primitive.ObjectID) ([]catalog.Showtime, error) {
	return nil, nil
}
func (fakeST) ListAdminShowtimes(context.Context, catalog.AdminShowtimeFilter) ([]catalog.Showtime, error) {
	return nil, nil
}
func (fakeST) UpdateShowtime(context.Context, *catalog.Showtime) error { return nil }
func (fakeST) DeleteShowtime(context.Context, primitive.ObjectID) error { return nil }

type fakeSC struct{ s *catalog.Screen }

func (f fakeSC) InsertScreen(context.Context, *catalog.Screen) error { return nil }
func (f fakeSC) FindScreenByID(context.Context, primitive.ObjectID) (*catalog.Screen, error) {
	return f.s, nil
}
func (f fakeSC) ListScreensByCinema(context.Context, primitive.ObjectID) ([]catalog.Screen, error) {
	return nil, nil
}
func (fakeSC) ListScreens(context.Context, *primitive.ObjectID) ([]catalog.Screen, error) { return nil, nil }
func (fakeSC) UpdateScreen(context.Context, *catalog.Screen) error                         { return nil }
func (fakeSC) DeleteScreen(context.Context, primitive.ObjectID) error                     { return nil }

type fakeMV struct{ m *catalog.Movie }

func (f fakeMV) InsertMovie(context.Context, *catalog.Movie) error { return nil }
func (f fakeMV) FindMovieByID(context.Context, primitive.ObjectID) (*catalog.Movie, error) {
	return f.m, nil
}
func (f fakeMV) ListMoviesByStatus(context.Context, string) ([]catalog.Movie, error) { return nil, nil }
func (f fakeMV) ListComingSoonMovies(context.Context) ([]catalog.Movie, error)      { return nil, nil }
func (f fakeMV) ListNonArchivedMovies(context.Context) ([]catalog.Movie, error)     { return nil, nil }
func (fakeMV) ListMovies(context.Context) ([]catalog.Movie, error)                  { return nil, nil }
func (fakeMV) UpdateMovie(context.Context, *catalog.Movie) error                    { return nil }
func (fakeMV) DeleteMovie(context.Context, primitive.ObjectID) error                { return nil }

type fakeCN struct{ c *catalog.Cinema }

func (f fakeCN) InsertCinema(context.Context, *catalog.Cinema) error { return nil }
func (f fakeCN) FindCinemaByID(context.Context, primitive.ObjectID) (*catalog.Cinema, error) {
	return f.c, nil
}
func (f fakeCN) ListCinemas(context.Context) ([]catalog.Cinema, error) { return nil, nil }
func (fakeCN) UpdateCinema(context.Context, *catalog.Cinema) error      { return nil }
func (fakeCN) DeleteCinema(context.Context, primitive.ObjectID) error   { return nil }
