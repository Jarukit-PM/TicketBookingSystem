package booking

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

// MongoRepository stores bookings in MongoDB.
type MongoRepository struct {
	coll *mongo.Collection
}

// NewMongoRepository returns a booking repository backed by the given database.
func NewMongoRepository(db *mongo.Database) *MongoRepository {
	return &MongoRepository{coll: db.Collection(CollectionName)}
}

func (r *MongoRepository) Insert(ctx context.Context, b *Booking) error {
	if b.ConfirmedAt.IsZero() {
		b.ConfirmedAt = time.Now().UTC()
	}
	res, err := r.coll.InsertOne(ctx, b)
	if err != nil {
		return fmt.Errorf("insert booking: %w", err)
	}
	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		b.ID = oid
	}
	return nil
}

func (r *MongoRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*Booking, error) {
	var b Booking
	err := r.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&b)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, fmt.Errorf("find booking by id: %w", err)
	}
	return &b, nil
}

func (r *MongoRepository) FindByBookingRef(ctx context.Context, ref string) (*Booking, error) {
	var b Booking
	err := r.coll.FindOne(ctx, bson.M{"bookingRef": ref}).Decode(&b)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, fmt.Errorf("find booking by ref: %w", err)
	}
	return &b, nil
}

func (r *MongoRepository) ListByUser(ctx context.Context, userID primitive.ObjectID) ([]Booking, error) {
	return r.ListConfirmedByUser(ctx, userID)
}

func (r *MongoRepository) ListConfirmedByUser(ctx context.Context, userID primitive.ObjectID) ([]Booking, error) {
	filter := bson.M{
		"userId": userID,
		"status": StatusConfirmed,
	}
	cur, err := r.coll.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "confirmedAt", Value: -1}}))
	if err != nil {
		return nil, fmt.Errorf("list confirmed bookings by user: %w", err)
	}
	defer cur.Close(ctx)

	var out []Booking
	if err := cur.All(ctx, &out); err != nil {
		return nil, fmt.Errorf("decode bookings: %w", err)
	}
	return out, nil
}

func (r *MongoRepository) CountConfirmedBetween(ctx context.Context, from, to time.Time) (int, error) {
	filter := bson.M{
		"status": StatusConfirmed,
		"confirmedAt": bson.M{
			"$gte": from,
			"$lt":  to,
		},
	}
	count, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("count confirmed bookings: %w", err)
	}
	return int(count), nil
}

func (r *MongoRepository) CountConfirmed(ctx context.Context) (int64, error) {
	count, err := r.coll.CountDocuments(ctx, bson.M{"status": StatusConfirmed})
	if err != nil {
		return 0, fmt.Errorf("count confirmed bookings: %w", err)
	}
	return count, nil
}

func (r *MongoRepository) ListConfirmedPage(ctx context.Context, skip, limit int) ([]Booking, error) {
	if limit <= 0 {
		limit = 20
	}
	if skip < 0 {
		skip = 0
	}
	filter := bson.M{"status": StatusConfirmed}
	opts := options.Find().
		SetSort(bson.D{{Key: "confirmedAt", Value: -1}}).
		SetSkip(int64(skip)).
		SetLimit(int64(limit))

	cur, err := r.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("list confirmed bookings page: %w", err)
	}
	defer cur.Close(ctx)

	var out []Booking
	if err := cur.All(ctx, &out); err != nil {
		return nil, fmt.Errorf("decode bookings: %w", err)
	}
	return out, nil
}

func (r *MongoRepository) ListRecentConfirmed(ctx context.Context, limit int) ([]Booking, error) {
	if limit <= 0 {
		limit = 10
	}
	filter := bson.M{"status": StatusConfirmed}
	opts := options.Find().
		SetSort(bson.D{{Key: "confirmedAt", Value: -1}}).
		SetLimit(int64(limit))

	cur, err := r.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("list recent confirmed bookings: %w", err)
	}
	defer cur.Close(ctx)

	var out []Booking
	if err := cur.All(ctx, &out); err != nil {
		return nil, fmt.Errorf("decode bookings: %w", err)
	}
	return out, nil
}

func (r *MongoRepository) CountConfirmedFiltered(ctx context.Context, filter ConfirmedFilter) (int64, error) {
	count, err := r.coll.CountDocuments(ctx, confirmedFilterBSON(filter))
	if err != nil {
		return 0, fmt.Errorf("count filtered confirmed bookings: %w", err)
	}
	return count, nil
}

func (r *MongoRepository) ListConfirmedFiltered(
	ctx context.Context,
	filter ConfirmedFilter,
	skip, limit int,
) ([]Booking, error) {
	if limit <= 0 {
		limit = 20
	}
	if skip < 0 {
		skip = 0
	}
	opts := options.Find().
		SetSort(bson.D{{Key: "confirmedAt", Value: -1}}).
		SetSkip(int64(skip)).
		SetLimit(int64(limit))

	cur, err := r.coll.Find(ctx, confirmedFilterBSON(filter), opts)
	if err != nil {
		return nil, fmt.Errorf("list filtered confirmed bookings: %w", err)
	}
	defer cur.Close(ctx)

	var out []Booking
	if err := cur.All(ctx, &out); err != nil {
		return nil, fmt.Errorf("decode bookings: %w", err)
	}
	return out, nil
}

func (r *MongoRepository) UpdateTicketToken(ctx context.Context, id primitive.ObjectID, token string) error {
	_, err := r.coll.UpdateByID(ctx, id, bson.M{"$set": bson.M{"ticketToken": token}})
	if err != nil {
		return fmt.Errorf("update ticket token: %w", err)
	}
	return nil
}

func (r *MongoRepository) ListConfirmedByShowtime(ctx context.Context, showtimeID primitive.ObjectID) ([]Booking, error) {
	filter := bson.M{
		"showtimeId": showtimeID,
		"status":     StatusConfirmed,
	}
	cur, err := r.coll.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "confirmedAt", Value: 1}}))
	if err != nil {
		return nil, fmt.Errorf("list bookings by showtime: %w", err)
	}
	defer cur.Close(ctx)

	var out []Booking
	if err := cur.All(ctx, &out); err != nil {
		return nil, fmt.Errorf("decode bookings: %w", err)
	}
	return out, nil
}
