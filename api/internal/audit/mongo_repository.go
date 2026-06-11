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

func (r *mongoAuditRepo) ListAuditLogs(ctx context.Context, page LogPage) ([]AuditLog, error) {
	opts := options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}})
	if page.Limit > 0 {
		opts.SetLimit(page.Limit)
	}
	if page.Skip > 0 {
		opts.SetSkip(page.Skip)
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
	return r.ListEmailLogs(ctx, LogPage{}, &bookingID)
}

func (r *mongoEmailLogRepo) ListEmailLogs(ctx context.Context, page LogPage, bookingID *primitive.ObjectID) ([]EmailLog, error) {
	filter := bson.M{}
	if bookingID != nil && !bookingID.IsZero() {
		filter["bookingId"] = *bookingID
	}
	opts := options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}})
	if page.Limit > 0 {
		opts.SetLimit(page.Limit)
	}
	if page.Skip > 0 {
		opts.SetSkip(page.Skip)
	}
	cur, err := r.coll.Find(ctx, filter, opts)
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
