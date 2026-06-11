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
	cur, err := r.coll.Find(ctx, bson.M{"userId": userID}, options.Find().SetSort(bson.D{{Key: "confirmedAt", Value: -1}}))
	if err != nil {
		return nil, fmt.Errorf("list bookings by user: %w", err)
	}
	defer cur.Close(ctx)

	var out []Booking
	if err := cur.All(ctx, &out); err != nil {
		return nil, fmt.Errorf("decode bookings: %w", err)
	}
	return out, nil
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
