package email

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hibiken/asynq"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/audit"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/booking"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/catalog"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/tasks"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/user"
)

const (
	StatusSent   = "SENT"
	StatusFailed = "FAILED"
)

// CatalogReader loads showtime metadata for confirmation emails.
type CatalogReader struct {
	Showtimes catalog.ShowtimeRepository
	Screens   catalog.ScreenRepository
	Movies    catalog.MovieRepository
	Cinemas   catalog.CinemaRepository
}

// Service sends booking confirmation emails from asynq tasks.
type Service struct {
	bookings booking.Repository
	users    user.Repository
	catalog  CatalogReader
	logs     audit.EmailLogRepository
	sender   Sender
	appURL   string
}

// NewService wires the email worker service.
func NewService(
	bookings booking.Repository,
	users user.Repository,
	catalog CatalogReader,
	logs audit.EmailLogRepository,
	sender Sender,
	appURL string,
) *Service {
	return &Service{
		bookings: bookings,
		users:    users,
		catalog:  catalog,
		logs:     logs,
		sender:   sender,
		appURL:   appURL,
	}
}

// HandleEmailSend processes an email:send asynq task.
func (s *Service) HandleEmailSend(ctx context.Context, t *asynq.Task) error {
	var payload tasks.EmailSendPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}
	id, err := primitive.ObjectIDFromHex(payload.BookingID)
	if err != nil {
		return err
	}
	b, err := s.bookings.FindByID(ctx, id)
	if err != nil || b == nil || b.Status != booking.StatusConfirmed {
		return fmt.Errorf("booking not found")
	}
	u, err := s.users.FindByID(ctx, b.UserID)
	if err != nil || u == nil || u.Email == "" {
		return fmt.Errorf("recipient not found")
	}
	st, err := s.catalog.Showtimes.FindShowtimeByID(ctx, b.ShowtimeID)
	if err != nil || st == nil {
		return fmt.Errorf("showtime not found")
	}
	sc, err := s.catalog.Screens.FindScreenByID(ctx, st.ScreenID)
	if err != nil || sc == nil {
		return fmt.Errorf("screen not found")
	}
	mv, err := s.catalog.Movies.FindMovieByID(ctx, st.MovieID)
	if err != nil || mv == nil {
		return fmt.Errorf("movie not found")
	}
	cn, err := s.catalog.Cinemas.FindCinemaByID(ctx, sc.CinemaID)
	if err != nil || cn == nil {
		return fmt.Errorf("cinema not found")
	}

	ticketURL := booking.TicketURL(s.appURL, b.BookingRef, b.TicketToken)
	locale := booking.ParseLocale(b.Locale)
	html, text, err := renderConfirmation(locale, confirmationData{
		BookingRef: b.BookingRef,
		MovieTitle: mv.Title,
		CinemaName: cn.Name,
		ScreenName: sc.Name,
		StartsAt:   st.StartsAt.UTC().Format(time.RFC1123),
		Seats:      formatSeats(b.Seats),
		Total:      fmt.Sprintf("%d", b.Total),
		TicketURL:  ticketURL,
	})
	if err != nil {
		return err
	}

	providerID, sendErr := s.sender.Send(ctx, Message{
		To:       u.Email,
		Subject:  confirmationSubject(locale, mv.Title),
		HTMLBody: html,
		TextBody: text,
	})
	status := StatusSent
	if sendErr != nil {
		status = StatusFailed
		log.Printf("email:send failed booking=%s: %v", payload.BookingID, sendErr)
	}
	_ = s.logs.InsertEmailLog(ctx, &audit.EmailLog{
		BookingID:  b.ID,
		Type:       audit.EmailTypeConfirmation,
		To:         u.Email,
		ProviderID: providerID,
		Status:     status,
		CreatedAt:  time.Now().UTC(),
	})
	return sendErr
}
