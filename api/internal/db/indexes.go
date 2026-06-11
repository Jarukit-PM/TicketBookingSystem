package db

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/audit"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/booking"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/catalog"
	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/user"
)

type indexSpec struct {
	collection string
	models     []mongo.IndexModel
}

// EnsureIndexes creates all MVP indexes idempotently on server boot.
func EnsureIndexes(ctx context.Context, db *mongo.Database) error {
	specs := []indexSpec{
		{
			collection: user.CollectionName,
			models: []mongo.IndexModel{
				{
					Keys:    bson.D{{Key: "email", Value: 1}},
					Options: options.Index().SetUnique(true).SetName("uniq_email"),
				},
				{
					Keys:    bson.D{{Key: "googleId", Value: 1}},
					Options: options.Index().SetUnique(true).SetSparse(true).SetName("uniq_googleId_sparse"),
				},
			},
		},
		{
			collection: catalog.CollectionMovies,
			models: []mongo.IndexModel{
				{
					Keys:    bson.D{{Key: "status", Value: 1}},
					Options: options.Index().SetName("idx_status"),
				},
			},
		},
		{
			collection: catalog.CollectionScreens,
			models: []mongo.IndexModel{
				{
					Keys:    bson.D{{Key: "cinemaId", Value: 1}},
					Options: options.Index().SetName("idx_cinemaId"),
				},
			},
		},
		{
			collection: catalog.CollectionShowtimes,
			models: []mongo.IndexModel{
				{
					Keys: bson.D{
						{Key: "screenId", Value: 1},
						{Key: "startsAt", Value: 1},
					},
					Options: options.Index().SetName("idx_screenId_startsAt"),
				},
				{
					Keys:    bson.D{{Key: "movieId", Value: 1}},
					Options: options.Index().SetName("idx_movieId"),
				},
			},
		},
		{
			collection: booking.CollectionName,
			models: []mongo.IndexModel{
				{
					Keys:    bson.D{{Key: "bookingRef", Value: 1}},
					Options: options.Index().SetUnique(true).SetName("uniq_bookingRef"),
				},
				{
					Keys:    bson.D{{Key: "userId", Value: 1}},
					Options: options.Index().SetName("idx_userId"),
				},
				{
					Keys: bson.D{
						{Key: "userId", Value: 1},
						{Key: "showtimeId", Value: 1},
					},
					Options: options.Index().SetName("idx_userId_showtimeId"),
				},
				{
					Keys: bson.D{
						{Key: "status", Value: 1},
						{Key: "confirmedAt", Value: -1},
					},
					Options: options.Index().SetName("idx_status_confirmedAt"),
				},
			},
		},
		{
			collection: audit.CollectionAuditLogs,
			models: []mongo.IndexModel{
				{
					Keys:    bson.D{{Key: "createdAt", Value: -1}},
					Options: options.Index().SetName("idx_createdAt"),
				},
			},
		},
		{
			collection: audit.CollectionEmailLogs,
			models: []mongo.IndexModel{
				{
					Keys:    bson.D{{Key: "bookingId", Value: 1}},
					Options: options.Index().SetName("idx_bookingId"),
				},
			},
		},
	}

	for _, spec := range specs {
		if _, err := db.Collection(spec.collection).Indexes().CreateMany(ctx, spec.models); err != nil {
			return fmt.Errorf("ensure indexes on %s: %w", spec.collection, err)
		}
	}
	return nil
}
