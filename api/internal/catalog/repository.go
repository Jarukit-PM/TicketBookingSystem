package catalog

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AdminShowtimeFilter narrows admin showtime listings.
type AdminShowtimeFilter struct {
	CinemaID *primitive.ObjectID
	MovieID  *primitive.ObjectID
	ScreenID *primitive.ObjectID
	From     *time.Time
	To       *time.Time
}

// CinemaRepository persists cinemas.
type CinemaRepository interface {
	InsertCinema(ctx context.Context, cinema *Cinema) error
	FindCinemaByID(ctx context.Context, id primitive.ObjectID) (*Cinema, error)
	ListCinemas(ctx context.Context) ([]Cinema, error)
	UpdateCinema(ctx context.Context, cinema *Cinema) error
	DeleteCinema(ctx context.Context, id primitive.ObjectID) error
}

// ScreenRepository persists screens.
type ScreenRepository interface {
	InsertScreen(ctx context.Context, screen *Screen) error
	FindScreenByID(ctx context.Context, id primitive.ObjectID) (*Screen, error)
	ListScreensByCinema(ctx context.Context, cinemaID primitive.ObjectID) ([]Screen, error)
	ListScreens(ctx context.Context, cinemaID *primitive.ObjectID) ([]Screen, error)
	UpdateScreen(ctx context.Context, screen *Screen) error
	DeleteScreen(ctx context.Context, id primitive.ObjectID) error
}

// MovieRepository persists movies.
type MovieRepository interface {
	InsertMovie(ctx context.Context, movie *Movie) error
	FindMovieByID(ctx context.Context, id primitive.ObjectID) (*Movie, error)
	ListMoviesByStatus(ctx context.Context, status string) ([]Movie, error)
	ListMovies(ctx context.Context) ([]Movie, error)
	ListComingSoonMovies(ctx context.Context) ([]Movie, error)
	ListNonArchivedMovies(ctx context.Context) ([]Movie, error)
	UpdateMovie(ctx context.Context, movie *Movie) error
	DeleteMovie(ctx context.Context, id primitive.ObjectID) error
}

// ShowtimeRepository persists showtimes.
type ShowtimeRepository interface {
	InsertShowtime(ctx context.Context, showtime *Showtime) error
	FindShowtimeByID(ctx context.Context, id primitive.ObjectID) (*Showtime, error)
	ListShowtimesByScreen(ctx context.Context, screenID primitive.ObjectID, from time.Time) ([]Showtime, error)
	ListShowtimesByMovie(ctx context.Context, movieID primitive.ObjectID) ([]Showtime, error)
	ListAdminShowtimes(ctx context.Context, filter AdminShowtimeFilter) ([]Showtime, error)
	ListShowtimesByScreens(ctx context.Context, screenIDs []primitive.ObjectID) ([]Showtime, error)
	ListShowtimesByCinemaMovie(ctx context.Context, screenIDs []primitive.ObjectID, movieID primitive.ObjectID) ([]Showtime, error)
	UpdateShowtime(ctx context.Context, showtime *Showtime) error
	DeleteShowtime(ctx context.Context, id primitive.ObjectID) error
}
