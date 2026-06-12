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
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/audit"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/auth"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/handler"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/user"
)

type memAuditLogs struct {
	logs []audit.AuditLog
}

func (r *memAuditLogs) InsertAuditLog(_ context.Context, log *audit.AuditLog) error {
	r.logs = append(r.logs, *log)
	return nil
}
func (r *memAuditLogs) ListAuditLogs(_ context.Context, page audit.LogPage, filter audit.AuditLogFilter) ([]audit.AuditLog, error) {
	var filtered []audit.AuditLog
	for _, log := range r.logs {
		if filter.Matches(log) {
			filtered = append(filtered, log)
		}
	}
	start := int(page.Skip)
	end := start + int(page.Limit)
	if start > len(filtered) {
		return []audit.AuditLog{}, nil
	}
	if end > len(filtered) {
		end = len(filtered)
	}
	return filtered[start:end], nil
}

type memEmailLogs struct {
	logs []audit.EmailLog
}

func (r *memEmailLogs) InsertEmailLog(_ context.Context, log *audit.EmailLog) error {
	r.logs = append(r.logs, *log)
	return nil
}
func (r *memEmailLogs) ListByBooking(_ context.Context, bookingID primitive.ObjectID) ([]audit.EmailLog, error) {
	var out []audit.EmailLog
	for _, l := range r.logs {
		if l.BookingID == bookingID {
			out = append(out, l)
		}
	}
	return out, nil
}
func (r *memEmailLogs) ListEmailLogs(_ context.Context, page audit.LogPage, filter audit.EmailLogFilter) ([]audit.EmailLog, error) {
	var filtered []audit.EmailLog
	for _, log := range r.logs {
		if filter.Matches(log) {
			filtered = append(filtered, log)
		}
	}
	start := int(page.Skip)
	end := start + int(page.Limit)
	if start > len(filtered) {
		return []audit.EmailLog{}, nil
	}
	if end > len(filtered) {
		end = len(filtered)
	}
	return filtered[start:end], nil
}

func setupAdminLogsRouter(t *testing.T, role string, auditRepo *memAuditLogs, emailRepo *memEmailLogs) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)

	svc := &admin.LogsService{AuditLogs: auditRepo, EmailLogs: emailRepo}
	deps := handler.AdminLogsDeps{Service: svc}

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
	adminGroup.GET("/audit-logs", handler.ListAdminAuditLogs(deps))
	adminGroup.GET("/email-logs", handler.ListAdminEmailLogs(deps))

	return r
}

func TestAdminAuditLogsRejectsCustomer(t *testing.T) {
	r := setupAdminLogsRouter(t, user.RoleCustomer, &memAuditLogs{}, &memEmailLogs{})

	req := httptest.NewRequest(http.MethodGet, "/api/admin/audit-logs", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403", w.Code)
	}
}

func TestAdminAuditLogsListsNewestFirst(t *testing.T) {
	auditRepo := &memAuditLogs{
		logs: []audit.AuditLog{
			{Action: "UPDATE", Entity: "movie", CreatedAt: time.Now().UTC()},
			{Action: "CREATE", Entity: "movie", CreatedAt: time.Now().UTC().Add(-time.Hour)},
		},
	}
	r := setupAdminLogsRouter(t, user.RoleAdmin, auditRepo, &memEmailLogs{})

	req := httptest.NewRequest(http.MethodGet, "/api/admin/audit-logs?limit=1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body %s", w.Code, w.Body.String())
	}
	body := w.Body.String()
	if !strings.Contains(body, "UPDATE") {
		t.Fatalf("body = %s, want newest audit action", body)
	}
}

func TestAdminAuditLogsFilterByAction(t *testing.T) {
	auditRepo := &memAuditLogs{
		logs: []audit.AuditLog{
			{Action: audit.ActionBookingSuccess, Entity: "booking"},
			{Action: audit.ActionCreate, Entity: "movie"},
		},
	}
	r := setupAdminLogsRouter(t, user.RoleAdmin, auditRepo, &memEmailLogs{})

	req := httptest.NewRequest(http.MethodGet, "/api/admin/audit-logs?action="+audit.ActionBookingSuccess, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body %s", w.Code, w.Body.String())
	}
	body := w.Body.String()
	if !strings.Contains(body, `"action":"`+audit.ActionBookingSuccess+`"`) {
		t.Fatalf("body = %s, want booking_success audit log", body)
	}
	if strings.Contains(body, `"action":"`+audit.ActionCreate+`"`) {
		t.Fatalf("body = %s, want create action filtered out", body)
	}
}

func TestAdminEmailLogsFilterByStatus(t *testing.T) {
	bookingID := primitive.NewObjectID()
	emailRepo := &memEmailLogs{
		logs: []audit.EmailLog{
			{BookingID: bookingID, Type: audit.EmailTypeConfirmation, Status: "SENT", To: "a@example.com"},
			{BookingID: bookingID, Type: audit.EmailTypeConfirmation, Status: "FAILED", To: "a@example.com"},
		},
	}
	r := setupAdminLogsRouter(t, user.RoleAdmin, &memAuditLogs{}, emailRepo)

	req := httptest.NewRequest(http.MethodGet, "/api/admin/email-logs?status=FAILED", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body %s", w.Code, w.Body.String())
	}
	body := w.Body.String()
	if !strings.Contains(body, `"status":"FAILED"`) {
		t.Fatalf("body = %s, want failed email log", body)
	}
	if strings.Contains(body, `"status":"SENT"`) {
		t.Fatalf("body = %s, want sent status filtered out", body)
	}
}

func TestAdminEmailLogsFilterByBooking(t *testing.T) {
	bookingID := primitive.NewObjectID()
	otherID := primitive.NewObjectID()
	emailRepo := &memEmailLogs{
		logs: []audit.EmailLog{
			{BookingID: bookingID, Type: audit.EmailTypeConfirmation, Status: "SENT", To: "a@example.com"},
			{BookingID: otherID, Type: audit.EmailTypeConfirmation, Status: "SENT", To: "b@example.com"},
		},
	}
	r := setupAdminLogsRouter(t, user.RoleAdmin, &memAuditLogs{}, emailRepo)

	req := httptest.NewRequest(http.MethodGet, "/api/admin/email-logs?bookingId="+bookingID.Hex(), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body %s", w.Code, w.Body.String())
	}
	body := w.Body.String()
	if !strings.Contains(body, "a@example.com") || strings.Contains(body, "b@example.com") {
		t.Fatalf("body = %s, want only matching booking email log", body)
	}
}
