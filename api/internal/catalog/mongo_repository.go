package catalog

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoRepositories groups catalog mongo repositories.
type MongoRepositories struct {
	Cinemas   CinemaRepository
	Screens   ScreenRepository
	Movies    MovieRepository
	Showtimes ShowtimeRepository
}

// NewMongoRepositories returns catalog repositories backed by the given database.
func NewMongoRepositories(db *mongo.Database) MongoRepositories {
	return MongoRepositories{
		Cinemas:   &mongoCinemaRepo{coll: db.Collection(CollectionCinemas)},
		Screens:   &mongoScreenRepo{coll: db.Collection(CollectionScreens)},
		Movies:    &mongoMovieRepo{coll: db.Collection(CollectionMovies)},
		Showtimes: &mongoShowtimeRepo{coll: db.Collection(CollectionShowtimes)},
	}
}

type mongoCinemaRepo struct {
	coll *mongo.Collection
}

func (r *mongoCinemaRepo) InsertCinema(ctx context.Context, cinema *Cinema) error {
	res, err := r.coll.InsertOne(ctx, cinema)
	if err != nil {
		return fmt.Errorf("insert cinema: %w", err)
	}
	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		cinema.ID = oid
	}
	return nil
}

func (r *mongoCinemaRepo) FindCinemaByID(ctx context.Context, id primitive.ObjectID) (*Cinema, error) {
	var c Cinema
	err := r.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&c)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, fmt.Errorf("find cinema: %w", err)
	}
	return &c, nil
}

func (r *mongoCinemaRepo) ListCinemas(ctx context.Context) ([]Cinema, error) {
	cur, err := r.coll.Find(ctx, bson.M{}, options.Find().SetSort(bson.D{{Key: "name", Value: 1}}))
	if err != nil {
		return nil, fmt.Errorf("list cinemas: %w", err)
	}
	defer cur.Close(ctx)

	var out []Cinema
	if err := cur.All(ctx, &out); err != nil {
		return nil, fmt.Errorf("decode cinemas: %w", err)
	}
	return out, nil
}

type mongoScreenRepo struct {
	coll *mongo.Collection
}

func (r *mongoScreenRepo) InsertScreen(ctx context.Context, screen *Screen) error {
	res, err := r.coll.InsertOne(ctx, screen)
	if err != nil {
		return fmt.Errorf("insert screen: %w", err)
	}
	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		screen.ID = oid
	}
	return nil
}

func (r *mongoScreenRepo) FindScreenByID(ctx context.Context, id primitive.ObjectID) (*Screen, error) {
	var s Screen
	err := r.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&s)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, fmt.Errorf("find screen: %w", err)
	}
	return &s, nil
}

func (r *mongoScreenRepo) ListScreensByCinema(ctx context.Context, cinemaID primitive.ObjectID) ([]Screen, error) {
	cur, err := r.coll.Find(ctx, bson.M{"cinemaId": cinemaID}, options.Find().SetSort(bson.D{{Key: "name", Value: 1}}))
	if err != nil {
		return nil, fmt.Errorf("list screens: %w", err)
	}
	defer cur.Close(ctx)

	var out []Screen
	if err := cur.All(ctx, &out); err != nil {
		return nil, fmt.Errorf("decode screens: %w", err)
	}
	return out, nil
}

type mongoMovieRepo struct {
	coll *mongo.Collection
}

func (r *mongoMovieRepo) InsertMovie(ctx context.Context, movie *Movie) error {
	res, err := r.coll.InsertOne(ctx, movie)
	if err != nil {
		return fmt.Errorf("insert movie: %w", err)
	}
	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		movie.ID = oid
	}
	return nil
}

func (r *mongoMovieRepo) FindMovieByID(ctx context.Context, id primitive.ObjectID) (*Movie, error) {
	var m Movie
	err := r.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&m)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, fmt.Errorf("find movie: %w", err)
	}
	return &m, nil
}

func (r *mongoMovieRepo) ListMoviesByStatus(ctx context.Context, status string) ([]Movie, error) {
	cur, err := r.coll.Find(ctx, bson.M{"status": status}, options.Find().SetSort(bson.D{{Key: "title", Value: 1}}))
	if err != nil {
		return nil, fmt.Errorf("list movies: %w", err)
	}
	defer cur.Close(ctx)

	var out []Movie
	if err := cur.All(ctx, &out); err != nil {
		return nil, fmt.Errorf("decode movies: %w", err)
	}
	return out, nil
}

type mongoShowtimeRepo struct {
	coll *mongo.Collection
}

func (r *mongoShowtimeRepo) InsertShowtime(ctx context.Context, showtime *Showtime) error {
	res, err := r.coll.InsertOne(ctx, showtime)
	if err != nil {
		return fmt.Errorf("insert showtime: %w", err)
	}
	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		showtime.ID = oid
	}
	return nil
}

func (r *mongoShowtimeRepo) FindShowtimeByID(ctx context.Context, id primitive.ObjectID) (*Showtime, error) {
	var s Showtime
	err := r.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&s)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, fmt.Errorf("find showtime: %w", err)
	}
	return &s, nil
}

func (r *mongoShowtimeRepo) ListShowtimesByScreen(ctx context.Context, screenID primitive.ObjectID, from time.Time) ([]Showtime, error) {
	filter := bson.M{
		"screenId": screenID,
		"startsAt": bson.M{"$gte": from},
	}
	cur, err := r.coll.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "startsAt", Value: 1}}))
	if err != nil {
		return nil, fmt.Errorf("list showtimes by screen: %w", err)
	}
	defer cur.Close(ctx)

	var out []Showtime
	if err := cur.All(ctx, &out); err != nil {
		return nil, fmt.Errorf("decode showtimes: %w", err)
	}
	return out, nil
}

func (r *mongoShowtimeRepo) ListShowtimesByMovie(ctx context.Context, movieID primitive.ObjectID) ([]Showtime, error) {
	cur, err := r.coll.Find(ctx, bson.M{"movieId": movieID}, options.Find().SetSort(bson.D{{Key: "startsAt", Value: 1}}))
	if err != nil {
		return nil, fmt.Errorf("list showtimes by movie: %w", err)
	}
	defer cur.Close(ctx)

	var out []Showtime
	if err := cur.All(ctx, &out); err != nil {
		return nil, fmt.Errorf("decode showtimes: %w", err)
	}
	return out, nil
}
