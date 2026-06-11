package user

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Repository persists and loads users.
type Repository interface {
	Insert(ctx context.Context, user *User) error
	FindByID(ctx context.Context, id primitive.ObjectID) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByGoogleID(ctx context.Context, googleID string) (*User, error)
	Update(ctx context.Context, user *User) error
}
