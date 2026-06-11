package admin

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/booking"
)

// ErrInvalidTicket is returned when a ticket ref/token pair cannot be resolved.
var ErrInvalidTicket = errors.New("invalid ticket")

// TicketResolveResult is returned by GET /api/admin/tickets/resolve.
type TicketResolveResult struct {
	UserID    string `json:"userId"`
	BookingID string `json:"bookingId"`
}

// TicketsService resolves scanned ticket QR codes for admin support lookup.
type TicketsService struct {
	Bookings     booking.Repository
	TicketSecret string
}

// Resolve validates ref + token and returns the booking owner for navigation.
func (s *TicketsService) Resolve(ctx context.Context, ref, token string) (*TicketResolveResult, error) {
	ref = strings.TrimSpace(ref)
	token = strings.TrimSpace(token)
	if ref == "" || token == "" {
		return nil, ErrInvalidTicket
	}

	b, err := s.Bookings.FindByBookingRef(ctx, ref)
	if err != nil {
		return nil, fmt.Errorf("find booking by ref: %w", err)
	}
	if b == nil || !booking.ValidateTicketToken(ref, token, b, s.TicketSecret) {
		return nil, ErrInvalidTicket
	}

	return &TicketResolveResult{
		UserID:    b.UserID.Hex(),
		BookingID: b.ID.Hex(),
	}, nil
}
