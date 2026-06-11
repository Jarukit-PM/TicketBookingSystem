package handler_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/admin"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/auth"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/booking"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/handler"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/user"
)

const testTicketSecret = "test-ticket-secret"

type adminTicketsRepo struct {
	byRef *booking.Booking
}

func (r *adminTicketsRepo) Insert(_ context.Context, _ *booking.Booking) error { return nil }
func (r *adminTicketsRepo) FindByID(_ context.Context, _ primitive.ObjectID) (*booking.Booking, error) {
	return nil, nil
}
func (r *adminTicketsRepo) FindByBookingRef(_ context.Context, ref string) (*booking.Booking, error) {
	if r.byRef != nil && r.byRef.BookingRef == ref {
		return r.byRef, nil
	}
	return nil, nil
}
func (r *adminTicketsRepo) ListByUser(ctx context.Context, userID primitive.ObjectID) ([]booking.Booking, error) {
	return r.ListConfirmedByUser(ctx, userID)
}
func (r *adminTicketsRepo) ListConfirmedByUser(_ context.Context, _ primitive.ObjectID) ([]booking.Booking, error) {
	return nil, nil
}
func (r *adminTicketsRepo) ListConfirmedByShowtime(_ context.Context, _ primitive.ObjectID) ([]booking.Booking, error) {
	return nil, nil
}
func (r *adminTicketsRepo) CountConfirmedBetween(_ context.Context, _, _ time.Time) (int, error) {
	return 0, nil
}
func (r *adminTicketsRepo) ListRecentConfirmed(_ context.Context, _ int) ([]booking.Booking, error) {
	return nil, nil
}
func (r *adminTicketsRepo) CountConfirmed(_ context.Context) (int64, error) { return 0, nil }
func (r *adminTicketsRepo) ListConfirmedPage(_ context.Context, _, _ int) ([]booking.Booking, error) {
	return nil, nil
}

func setupAdminTicketsRouter(t *testing.T, role string, repo *adminTicketsRepo) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)

	svc := &admin.TicketsService{
		Bookings:     repo,
		TicketSecret: testTicketSecret,
	}
	deps := handler.AdminTicketsDeps{Service: svc}

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
	adminGroup.GET("/tickets/resolve", handler.ResolveAdminTicket(deps))

	return r
}

func TestAdminTicketsResolveRejectsCustomer(t *testing.T) {
	r := setupAdminTicketsRouter(t, user.RoleCustomer, &adminTicketsRepo{})

	req := httptest.NewRequest(http.MethodGet, "/api/admin/tickets/resolve?ref=TBS-ABC&t=bad", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403", w.Code)
	}
}

func TestAdminTicketsResolveValidToken(t *testing.T) {
	ref := "TBS-VALID1"
	bookingID := primitive.NewObjectID()
	userID := primitive.NewObjectID()
	ticketToken := booking.SignTicketToken(testTicketSecret, ref, bookingID.Hex())

	repo := &adminTicketsRepo{
		byRef: &booking.Booking{
			ID:          bookingID,
			UserID:      userID,
			BookingRef:  ref,
			TicketToken: ticketToken,
			Status:      booking.StatusConfirmed,
		},
	}
	r := setupAdminTicketsRouter(t, user.RoleAdmin, repo)

	req := httptest.NewRequest(http.MethodGet, "/api/admin/tickets/resolve?ref="+ref+"&t="+ticketToken, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body %s", w.Code, w.Body.String())
	}
	body := w.Body.String()
	if !strings.Contains(body, userID.Hex()) || !strings.Contains(body, bookingID.Hex()) {
		t.Fatalf("body = %s, want userId and bookingId", body)
	}
}

func TestAdminTicketsResolveRejectsTamperedToken(t *testing.T) {
	ref := "TBS-TAMPER"
	bookingID := primitive.NewObjectID()
	ticketToken := booking.SignTicketToken(testTicketSecret, ref, bookingID.Hex())

	repo := &adminTicketsRepo{
		byRef: &booking.Booking{
			ID:          bookingID,
			UserID:      primitive.NewObjectID(),
			BookingRef:  ref,
			TicketToken: ticketToken,
			Status:      booking.StatusConfirmed,
		},
	}
	r := setupAdminTicketsRouter(t, user.RoleAdmin, repo)

	req := httptest.NewRequest(http.MethodGet, "/api/admin/tickets/resolve?ref="+ref+"&t="+ticketToken+"x", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", w.Code)
	}
}
