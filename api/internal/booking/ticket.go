package booking

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/skip2/go-qrcode"
	"go.mongodb.org/mongo-driver/bson/primitive"

)

// TicketDetail is returned by GET /api/bookings/:id/ticket.
type TicketDetail struct {
	BookingRef  string   `json:"bookingRef"`
	TicketURL   string   `json:"ticketUrl"`
	QRPngBase64 string   `json:"qrPngBase64"`
	Seats       []string `json:"seats"`
	Total       int64    `json:"total"`
	MovieTitle  string   `json:"movieTitle"`
	CinemaName  string   `json:"cinemaName"`
	ScreenName  string   `json:"screenName"`
	StartsAt    string   `json:"startsAt"`
}

// GetTicket returns ticket metadata and a QR PNG for the booking owner.
func (s *Service) GetTicket(ctx context.Context, userID, bookingID string) (*TicketDetail, error) {
	id, err := primitive.ObjectIDFromHex(bookingID)
	if err != nil {
		return nil, ErrBookingNotFound
	}
	userOID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, ErrForbidden
	}

	b, err := s.bookings.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if b == nil || b.Status != StatusConfirmed {
		return nil, ErrBookingNotFound
	}
	if b.UserID != userOID {
		return nil, ErrForbidden
	}

	showtime, err := s.showtimes.FindShowtimeByID(ctx, b.ShowtimeID)
	if err != nil || showtime == nil {
		return nil, fmt.Errorf("showtime not found")
	}
	screen, err := s.screens.FindScreenByID(ctx, showtime.ScreenID)
	if err != nil || screen == nil {
		return nil, fmt.Errorf("screen not found")
	}
	cinema, err := s.cinemas.FindCinemaByID(ctx, screen.CinemaID)
	if err != nil || cinema == nil {
		return nil, fmt.Errorf("cinema not found")
	}
	movie, err := s.movies.FindMovieByID(ctx, showtime.MovieID)
	if err != nil || movie == nil {
		return nil, fmt.Errorf("movie not found")
	}

	ticketURL := TicketURL(s.appURL, b.BookingRef, b.TicketToken)
	png, err := qrcode.Encode(ticketURL, qrcode.Medium, 256)
	if err != nil {
		return nil, fmt.Errorf("encode qr: %w", err)
	}

	return &TicketDetail{
		BookingRef:  b.BookingRef,
		TicketURL:   ticketURL,
		QRPngBase64: base64.StdEncoding.EncodeToString(png),
		Seats:       b.Seats,
		Total:       b.Total,
		MovieTitle:  movie.Title,
		CinemaName:  cinema.Name,
		ScreenName:  screen.Name,
		StartsAt:    showtime.StartsAt.UTC().Format("2006-01-02T15:04:05Z"),
	}, nil
}
