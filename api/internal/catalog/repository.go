package catalog

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CinemaRepository persists cinemas.
type CinemaRepository interface {
	InsertCinema(ctx context.Context, cinema *Cinema) error
	FindCinemaByID(ctx context.Context, id primitive.ObjectID) (*Cinema, error)
	ListCinemas(ctx context.Context) ([]Cinema, error)
}

// ScreenRepository persists screens.
type ScreenRepository interface {
	InsertScreen(ctx context.Context, screen *Screen) error
	FindScreenByID(ctx context.Context, id primitive.ObjectID) (*Screen, error)
	ListScreensByCinema(ctx context.Context, cinemaID primitive.ObjectID) ([]Screen, error)
}

// MovieRepository persists movies.
type MovieRepository interface {
	InsertMovie(ctx context.Context, movie *Movie) error
	FindMovieByID(ctx context.Context, id primitive.ObjectID) (*Movie, error)
	ListMoviesByStatus(ctx context.Context, status string) ([]Movie, error)
}

// ShowtimeRepository persists showtimes.
type ShowtimeRepository interface {
	InsertShowtime(ctx context.Context, showtime *Showtime) error
	FindShowtimeByID(ctx context.Context, id primitive.ObjectID) (*Showtime, error)
	ListShowtimesByScreen(ctx context.Context, screenID primitive.ObjectID, from time.Time) ([]Showtime, error)
	ListShowtimesByMovie(ctx context.Context, movieID primitive.ObjectID) ([]Showtime, error)
}
