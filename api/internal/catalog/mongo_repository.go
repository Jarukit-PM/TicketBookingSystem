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

func (r *mongoCinemaRepo) UpdateCinema(ctx context.Context, cinema *Cinema) error {
	_, err := r.coll.ReplaceOne(ctx, bson.M{"_id": cinema.ID}, cinema)
	if err != nil {
		return fmt.Errorf("update cinema: %w", err)
	}
	return nil
}

func (r *mongoCinemaRepo) DeleteCinema(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.coll.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("delete cinema: %w", err)
	}
	return nil
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
	return r.ListScreens(ctx, &cinemaID)
}

func (r *mongoScreenRepo) ListScreens(ctx context.Context, cinemaID *primitive.ObjectID) ([]Screen, error) {
	filter := bson.M{}
	if cinemaID != nil {
		filter["cinemaId"] = *cinemaID
	}
	cur, err := r.coll.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "name", Value: 1}}))
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

func (r *mongoScreenRepo) UpdateScreen(ctx context.Context, screen *Screen) error {
	_, err := r.coll.ReplaceOne(ctx, bson.M{"_id": screen.ID}, screen)
	if err != nil {
		return fmt.Errorf("update screen: %w", err)
	}
	return nil
}

func (r *mongoScreenRepo) DeleteScreen(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.coll.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("delete screen: %w", err)
	}
	return nil
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

func (r *mongoMovieRepo) ListMovies(ctx context.Context) ([]Movie, error) {
	cur, err := r.coll.Find(ctx, bson.M{}, options.Find().SetSort(bson.D{{Key: "title", Value: 1}}))
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

func (r *mongoMovieRepo) UpdateMovie(ctx context.Context, movie *Movie) error {
	_, err := r.coll.ReplaceOne(ctx, bson.M{"_id": movie.ID}, movie)
	if err != nil {
		return fmt.Errorf("update movie: %w", err)
	}
	return nil
}

func (r *mongoMovieRepo) DeleteMovie(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.coll.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("delete movie: %w", err)
	}
	return nil
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

func (r *mongoShowtimeRepo) ListAdminShowtimes(ctx context.Context, filter AdminShowtimeFilter) ([]Showtime, error) {
	query := bson.M{}
	if filter.MovieID != nil {
		query["movieId"] = *filter.MovieID
	}
	if filter.ScreenID != nil {
		query["screenId"] = *filter.ScreenID
	}
	if filter.From != nil || filter.To != nil {
		rangeFilter := bson.M{}
		if filter.From != nil {
			rangeFilter["$gte"] = *filter.From
		}
		if filter.To != nil {
			rangeFilter["$lte"] = *filter.To
		}
		query["startsAt"] = rangeFilter
	}

	cur, err := r.coll.Find(ctx, query, options.Find().SetSort(bson.D{{Key: "startsAt", Value: 1}}))
	if err != nil {
		return nil, fmt.Errorf("list admin showtimes: %w", err)
	}
	defer cur.Close(ctx)

	var out []Showtime
	if err := cur.All(ctx, &out); err != nil {
		return nil, fmt.Errorf("decode showtimes: %w", err)
	}

	if filter.CinemaID == nil {
		return out, nil
	}

	screenRepo := &mongoScreenRepo{coll: r.coll.Database().Collection(CollectionScreens)}
	screens, err := screenRepo.ListScreensByCinema(ctx, *filter.CinemaID)
	if err != nil {
		return nil, err
	}
	screenIDs := make(map[primitive.ObjectID]struct{}, len(screens))
	for _, s := range screens {
		screenIDs[s.ID] = struct{}{}
	}

	filtered := make([]Showtime, 0, len(out))
	for _, st := range out {
		if _, ok := screenIDs[st.ScreenID]; ok {
			filtered = append(filtered, st)
		}
	}
	return filtered, nil
}

func (r *mongoShowtimeRepo) UpdateShowtime(ctx context.Context, showtime *Showtime) error {
	_, err := r.coll.ReplaceOne(ctx, bson.M{"_id": showtime.ID}, showtime)
	if err != nil {
		return fmt.Errorf("update showtime: %w", err)
	}
	return nil
}

func (r *mongoShowtimeRepo) DeleteShowtime(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.coll.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("delete showtime: %w", err)
	}
	return nil
}
