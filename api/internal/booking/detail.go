package booking

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/user"
)

type ListMovie struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	PosterURL string `json:"posterUrl"`
}

type ListCinema struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ListScreen struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ListItem struct {
	ID          string     `json:"id"`
	BookingRef  string     `json:"bookingRef"`
	ShowtimeID  string     `json:"showtimeId"`
	Seats       []string   `json:"seats"`
	Total       int64      `json:"total"`
	Status      string     `json:"status"`
	Locale      string     `json:"locale"`
	ConfirmedAt time.Time  `json:"confirmedAt"`
	StartsAt    time.Time  `json:"startsAt"`
	Movie       ListMovie  `json:"movie"`
	Cinema      ListCinema `json:"cinema"`
	Screen      ListScreen `json:"screen"`
}

func (s *Service) ListMine(ctx context.Context, userID string, upcoming bool) ([]ListItem, error) {
	userOID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user id: %w", err)
	}
	bookings, err := s.bookings.ListConfirmedByUser(ctx, userOID)
	if err != nil {
		return nil, err
	}
	items := make([]ListItem, 0, len(bookings))
	for _, b := range bookings {
		item, cinemaTZ, err := s.enrichBooking(ctx, b)
		if err != nil {
			return nil, err
		}
		isUpcoming, err := showtimeIsUpcoming(item.StartsAt, cinemaTZ, s.now())
		if err != nil {
			return nil, err
		}
		if upcoming != isUpcoming {
			continue
		}
		items = append(items, item)
	}
	return items, nil
}

func (s *Service) GetDetail(ctx context.Context, userID, role, bookingID string) (*ListItem, error) {
	id, err := primitive.ObjectIDFromHex(bookingID)
	if err != nil {
		return nil, ErrBookingNotFound
	}
	b, err := s.bookings.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if b == nil || b.Status != StatusConfirmed {
		return nil, ErrBookingNotFound
	}
	if role != user.RoleAdmin {
		userOID, err := primitive.ObjectIDFromHex(userID)
		if err != nil || b.UserID != userOID {
			return nil, ErrForbidden
		}
	}
	item, _, err := s.enrichBooking(ctx, *b)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (s *Service) enrichBooking(ctx context.Context, b Booking) (ListItem, string, error) {
	showtime, err := s.showtimes.FindShowtimeByID(ctx, b.ShowtimeID)
	if err != nil || showtime == nil {
		return ListItem{}, "", fmt.Errorf("showtime not found")
	}
	screen, err := s.screens.FindScreenByID(ctx, showtime.ScreenID)
	if err != nil || screen == nil {
		return ListItem{}, "", fmt.Errorf("screen not found")
	}
	cinema, err := s.cinemas.FindCinemaByID(ctx, screen.CinemaID)
	if err != nil || cinema == nil {
		return ListItem{}, "", fmt.Errorf("cinema not found")
	}
	movie, err := s.movies.FindMovieByID(ctx, showtime.MovieID)
	if err != nil || movie == nil {
		return ListItem{}, "", fmt.Errorf("movie not found")
	}
	return ListItem{
		ID:          b.ID.Hex(),
		BookingRef:  b.BookingRef,
		ShowtimeID:  b.ShowtimeID.Hex(),
		Seats:       b.Seats,
		Total:       b.Total,
		Status:      b.Status,
		Locale:      ParseLocale(b.Locale),
		ConfirmedAt: b.ConfirmedAt,
		StartsAt:    showtime.StartsAt,
		Movie:       ListMovie{ID: movie.ID.Hex(), Title: movie.Title, PosterURL: movie.PosterURL},
		Cinema:      ListCinema{ID: cinema.ID.Hex(), Name: cinema.Name},
		Screen:      ListScreen{ID: screen.ID.Hex(), Name: screen.Name},
	}, cinema.Timezone, nil
}

func showtimeIsUpcoming(startsAt time.Time, cinemaTimezone string, now time.Time) (bool, error) {
	loc, err := time.LoadLocation(cinemaTimezone)
	if err != nil {
		return false, fmt.Errorf("load timezone %q: %w", cinemaTimezone, err)
	}
	return startsAt.In(loc).After(now.In(loc)), nil
}
