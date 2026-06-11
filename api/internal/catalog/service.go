package catalog

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrCinemaNotFound = errors.New("cinema not found")
	ErrMovieNotFound  = errors.New("movie not found")
)

// ShowtimeView is a showtime enriched with screen metadata for API responses.
type ShowtimeView struct {
	ID         primitive.ObjectID `json:"id"`
	MovieID    primitive.ObjectID `json:"movieId"`
	ScreenID   primitive.ObjectID `json:"screenId"`
	ScreenName string             `json:"screenName"`
	StartsAt   time.Time          `json:"startsAt"`
	PriceTiers PriceTiers         `json:"priceTiers"`
	Status     string             `json:"status"`
}

// MovieDetailView is a movie with upcoming showtimes at a cinema.
type MovieDetailView struct {
	Movie
	Showtimes []ShowtimeView `json:"showtimes"`
}

// Service implements public catalog browse queries.
type Service struct {
	Cinemas   CinemaRepository
	Screens   ScreenRepository
	Movies    MovieRepository
	Showtimes ShowtimeRepository
	Now       func() time.Time
}

// ListCinemas returns all cinemas sorted by name.
func (s *Service) ListCinemas(ctx context.Context) ([]Cinema, error) {
	return s.Cinemas.ListCinemas(ctx)
}

// BrowseMovies returns movies for the given cinema tab.
func (s *Service) BrowseMovies(ctx context.Context, cinemaID primitive.ObjectID, tab BrowseTab) ([]Movie, error) {
	cinema, err := s.requireCinema(ctx, cinemaID)
	if err != nil {
		return nil, err
	}

	switch tab {
	case BrowseTabNowShowing:
		return s.listNowShowing(ctx, cinema)
	case BrowseTabComingSoon:
		return s.Movies.ListComingSoonMovies(ctx)
	default:
		return nil, fmt.Errorf("unsupported tab %q", tab)
	}
}

func (s *Service) listNowShowing(ctx context.Context, cinema *Cinema) ([]Movie, error) {
	now, err := s.cinemaNow(cinema)
	if err != nil {
		return nil, err
	}

	screenIDs, err := s.screenIDsForCinema(ctx, cinema.ID)
	if err != nil {
		return nil, err
	}

	showtimes, err := s.Showtimes.ListShowtimesByScreens(ctx, screenIDs)
	if err != nil {
		return nil, err
	}

	futureIDs := FutureShowtimeMovieIDs(showtimes, now)
	movies, err := s.Movies.ListNonArchivedMovies(ctx)
	if err != nil {
		return nil, err
	}

	return FilterNowShowing(movies, futureIDs), nil
}

// GetMovieDetail returns a movie and its future showtimes at a cinema.
func (s *Service) GetMovieDetail(ctx context.Context, movieID, cinemaID primitive.ObjectID) (*MovieDetailView, error) {
	cinema, err := s.requireCinema(ctx, cinemaID)
	if err != nil {
		return nil, err
	}

	movie, err := s.Movies.FindMovieByID(ctx, movieID)
	if err != nil {
		return nil, err
	}
	if movie == nil || movie.Status == MovieStatusArchived {
		return nil, ErrMovieNotFound
	}

	views, err := s.ListShowtimes(ctx, cinemaID, movieID, nil)
	if err != nil {
		return nil, err
	}

	_ = cinema // cinema validated for timezone in ListShowtimes

	return &MovieDetailView{
		Movie:     *movie,
		Showtimes: views,
	}, nil
}

// ListShowtimes returns future showtimes for a movie at a cinema.
func (s *Service) ListShowtimes(
	ctx context.Context,
	cinemaID, movieID primitive.ObjectID,
	date *time.Time,
) ([]ShowtimeView, error) {
	cinema, err := s.requireCinema(ctx, cinemaID)
	if err != nil {
		return nil, err
	}

	now, err := s.cinemaNow(cinema)
	if err != nil {
		return nil, err
	}

	screenIDs, screensByID, err := s.screensForCinema(ctx, cinema.ID)
	if err != nil {
		return nil, err
	}

	showtimes, err := s.Showtimes.ListShowtimesByCinemaMovie(ctx, screenIDs, movieID)
	if err != nil {
		return nil, err
	}

	future := FilterFutureShowtimes(showtimes, now)
	if date != nil {
		future = filterShowtimesByLocalDate(future, cinema.Timezone, *date)
	}

	return enrichShowtimes(future, screensByID), nil
}

func (s *Service) requireCinema(ctx context.Context, cinemaID primitive.ObjectID) (*Cinema, error) {
	cinema, err := s.Cinemas.FindCinemaByID(ctx, cinemaID)
	if err != nil {
		return nil, err
	}
	if cinema == nil {
		return nil, ErrCinemaNotFound
	}
	return cinema, nil
}

func (s *Service) cinemaNow(cinema *Cinema) (time.Time, error) {
	if s.Now != nil {
		return s.Now().In(mustLoadLocation(cinema.Timezone)), nil
	}
	return NowInTimezone(cinema.Timezone)
}

func mustLoadLocation(tz string) *time.Location {
	loc, err := time.LoadLocation(tz)
	if err != nil {
		return time.UTC
	}
	return loc
}

func (s *Service) screenIDsForCinema(ctx context.Context, cinemaID primitive.ObjectID) ([]primitive.ObjectID, error) {
	screens, err := s.Screens.ListScreensByCinema(ctx, cinemaID)
	if err != nil {
		return nil, err
	}
	ids := make([]primitive.ObjectID, len(screens))
	for i, sc := range screens {
		ids[i] = sc.ID
	}
	return ids, nil
}

func (s *Service) screensForCinema(ctx context.Context, cinemaID primitive.ObjectID) ([]primitive.ObjectID, map[primitive.ObjectID]Screen, error) {
	screens, err := s.Screens.ListScreensByCinema(ctx, cinemaID)
	if err != nil {
		return nil, nil, err
	}
	ids := make([]primitive.ObjectID, len(screens))
	byID := make(map[primitive.ObjectID]Screen, len(screens))
	for i, sc := range screens {
		ids[i] = sc.ID
		byID[sc.ID] = sc
	}
	return ids, byID, nil
}

func enrichShowtimes(showtimes []Showtime, screensByID map[primitive.ObjectID]Screen) []ShowtimeView {
	out := make([]ShowtimeView, 0, len(showtimes))
	for _, st := range showtimes {
		screenName := ""
		if sc, ok := screensByID[st.ScreenID]; ok {
			screenName = sc.Name
		}
		out = append(out, ShowtimeView{
			ID:         st.ID,
			MovieID:    st.MovieID,
			ScreenID:   st.ScreenID,
			ScreenName: screenName,
			StartsAt:   st.StartsAt,
			PriceTiers: st.PriceTiers,
			Status:     st.Status,
		})
	}
	return out
}

func filterShowtimesByLocalDate(showtimes []Showtime, timezone string, date time.Time) []Showtime {
	loc := mustLoadLocation(timezone)
	dayStart := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, loc)
	dayEnd := dayStart.Add(24 * time.Hour)

	filtered := make([]Showtime, 0, len(showtimes))
	for _, st := range showtimes {
		local := st.StartsAt.In(loc)
		if !local.Before(dayStart) && local.Before(dayEnd) {
			filtered = append(filtered, st)
		}
	}
	return filtered
}
