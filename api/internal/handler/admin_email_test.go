package handler_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/auth"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/booking"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/handler"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/user"
)

type resendBookingsRepo struct {
	byID map[primitive.ObjectID]*booking.Booking
}

func (r *resendBookingsRepo) Insert(context.Context, *booking.Booking) error { return nil }
func (r *resendBookingsRepo) FindByID(_ context.Context, id primitive.ObjectID) (*booking.Booking, error) {
	if r.byID == nil {
		return nil, nil
	}
	return r.byID[id], nil
}
func (r *resendBookingsRepo) FindByBookingRef(context.Context, string) (*booking.Booking, error) {
	return nil, nil
}
func (r *resendBookingsRepo) ListByUser(context.Context, primitive.ObjectID) ([]booking.Booking, error) {
	return nil, nil
}
func (r *resendBookingsRepo) ListConfirmedByUser(context.Context, primitive.ObjectID) ([]booking.Booking, error) {
	return nil, nil
}
func (r *resendBookingsRepo) CountConfirmedBetween(context.Context, time.Time, time.Time) (int, error) {
	return 0, nil
}
func (r *resendBookingsRepo) ListRecentConfirmed(context.Context, int) ([]booking.Booking, error) {
	return nil, nil
}
func (r *resendBookingsRepo) CountConfirmed(context.Context) (int64, error) { return 0, nil }
func (r *resendBookingsRepo) CountConfirmedFiltered(ctx context.Context, _ booking.ConfirmedFilter) (int64, error) {
	return 0, nil
}
func (r *resendBookingsRepo) ListConfirmedPage(context.Context, int, int) ([]booking.Booking, error) {
	return nil, nil
}
func (r *resendBookingsRepo) ListConfirmedFiltered(ctx context.Context, _ booking.ConfirmedFilter, _, _ int) ([]booking.Booking, error) {
	return nil, nil
}
func (r *resendBookingsRepo) ListConfirmedByShowtime(context.Context, primitive.ObjectID) ([]booking.Booking, error) {
	return nil, nil
}
func (r *resendBookingsRepo) UpdateTicketToken(context.Context, primitive.ObjectID, string) error {
	return nil
}

type fakeEnqueuer struct {
	ids []string
}

func (f *fakeEnqueuer) EnqueueEmailSend(_ context.Context, bookingID string) error {
	f.ids = append(f.ids, bookingID)
	return nil
}

func setupAdminEmailRouter(t *testing.T, role string, bookings booking.Repository, tasks *fakeEnqueuer) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)

	deps := handler.AdminEmailDeps{Bookings: bookings, Tasks: tasks}
	tokens := auth.NewTokenService("test-secret", time.Hour)
	authSvc := auth.NewService(nil, tokens, auth.NewLoginRateLimiter(nil), "")
	authMw := auth.MiddlewareDeps{Service: authSvc}

	userID := primitive.NewObjectID()
	token, _, err := tokens.Issue(userID, role)
	if err != nil {
		t.Fatal(err)
	}

	r := gin.New()
	adminGroup := r.Group("/api/admin")
	adminGroup.Use(func(c *gin.Context) {
		c.Request.AddCookie(&http.Cookie{Name: auth.CookieName, Value: token})
		c.Next()
	})
	adminGroup.Use(auth.RequireAuth(authMw), auth.RequireAdmin(authMw))
	adminGroup.POST("/bookings/:id/resend-email", handler.ResendBookingEmail(deps))
	return r
}

func TestResendBookingEmailQueuesTask(t *testing.T) {
	bookingID := primitive.NewObjectID()
	repo := &resendBookingsRepo{
		byID: map[primitive.ObjectID]*booking.Booking{
			bookingID: {ID: bookingID, Status: booking.StatusConfirmed},
		},
	}
	enq := &fakeEnqueuer{}
	r := setupAdminEmailRouter(t, user.RoleAdmin, repo, enq)

	req := httptest.NewRequest(http.MethodPost, "/api/admin/bookings/"+bookingID.Hex()+"/resend-email", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", w.Code, w.Body.String())
	}
	if len(enq.ids) != 1 || enq.ids[0] != bookingID.Hex() {
		t.Fatalf("enqueued = %v", enq.ids)
	}
}

func TestResendBookingEmailNotFound(t *testing.T) {
	r := setupAdminEmailRouter(t, user.RoleAdmin, &resendBookingsRepo{byID: map[primitive.ObjectID]*booking.Booking{}}, &fakeEnqueuer{})

	req := httptest.NewRequest(http.MethodPost, "/api/admin/bookings/"+primitive.NewObjectID().Hex()+"/resend-email", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", w.Code)
	}
}
