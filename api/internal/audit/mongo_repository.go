package audit

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoRepositories groups audit mongo repositories.
type MongoRepositories struct {
	AuditLogs AuditRepository
	EmailLogs EmailLogRepository
}

// NewMongoRepositories returns audit repositories backed by the given database.
func NewMongoRepositories(db *mongo.Database) MongoRepositories {
	return MongoRepositories{
		AuditLogs: &mongoAuditRepo{coll: db.Collection(CollectionAuditLogs)},
		EmailLogs: &mongoEmailLogRepo{coll: db.Collection(CollectionEmailLogs)},
	}
}

type mongoAuditRepo struct {
	coll *mongo.Collection
}

func (r *mongoAuditRepo) InsertAuditLog(ctx context.Context, log *AuditLog) error {
	if log.CreatedAt.IsZero() {
		log.CreatedAt = time.Now().UTC()
	}
	res, err := r.coll.InsertOne(ctx, log)
	if err != nil {
		return fmt.Errorf("insert audit log: %w", err)
	}
	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		log.ID = oid
	}
	return nil
}

func (r *mongoAuditRepo) ListAuditLogs(ctx context.Context, limit int64) ([]AuditLog, error) {
	opts := options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}})
	if limit > 0 {
		opts.SetLimit(limit)
	}
	cur, err := r.coll.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, fmt.Errorf("list audit logs: %w", err)
	}
	defer cur.Close(ctx)

	var out []AuditLog
	if err := cur.All(ctx, &out); err != nil {
		return nil, fmt.Errorf("decode audit logs: %w", err)
	}
	return out, nil
}

type mongoEmailLogRepo struct {
	coll *mongo.Collection
}

func (r *mongoEmailLogRepo) InsertEmailLog(ctx context.Context, log *EmailLog) error {
	if log.CreatedAt.IsZero() {
		log.CreatedAt = time.Now().UTC()
	}
	res, err := r.coll.InsertOne(ctx, log)
	if err != nil {
		return fmt.Errorf("insert email log: %w", err)
	}
	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		log.ID = oid
	}
	return nil
}

func (r *mongoEmailLogRepo) ListByBooking(ctx context.Context, bookingID primitive.ObjectID) ([]EmailLog, error) {
	cur, err := r.coll.Find(ctx, bson.M{"bookingId": bookingID}, options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}}))
	if err != nil {
		return nil, fmt.Errorf("list email logs: %w", err)
	}
	defer cur.Close(ctx)

	var out []EmailLog
	if err := cur.All(ctx, &out); err != nil {
		return nil, fmt.Errorf("decode email logs: %w", err)
	}
	return out, nil
}
