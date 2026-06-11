package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// MongoRepository stores users in MongoDB.
type MongoRepository struct {
	coll *mongo.Collection
}

// NewMongoRepository returns a user repository backed by the given database.
func NewMongoRepository(db *mongo.Database) *MongoRepository {
	return &MongoRepository{coll: db.Collection(CollectionName)}
}

func (r *MongoRepository) Insert(ctx context.Context, user *User) error {
	if user.CreatedAt.IsZero() {
		user.CreatedAt = time.Now().UTC()
	}
	res, err := r.coll.InsertOne(ctx, user)
	if err != nil {
		return fmt.Errorf("insert user: %w", err)
	}
	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		user.ID = oid
	}
	return nil
}

func (r *MongoRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*User, error) {
	var u User
	err := r.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&u)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, fmt.Errorf("find user by id: %w", err)
	}
	return &u, nil
}

func (r *MongoRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
	var u User
	err := r.coll.FindOne(ctx, bson.M{"email": email}).Decode(&u)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, fmt.Errorf("find user by email: %w", err)
	}
	return &u, nil
}

func (r *MongoRepository) FindByGoogleID(ctx context.Context, googleID string) (*User, error) {
	var u User
	err := r.coll.FindOne(ctx, bson.M{"googleId": googleID}).Decode(&u)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, fmt.Errorf("find user by google id: %w", err)
	}
	return &u, nil
}

func (r *MongoRepository) Update(ctx context.Context, user *User) error {
	_, err := r.coll.ReplaceOne(ctx, bson.M{"_id": user.ID}, user)
	if err != nil {
		return fmt.Errorf("update user: %w", err)
	}
	return nil
}
