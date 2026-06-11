package user

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const CollectionName = "users"

const (
	RoleCustomer = "customer"
	RoleAdmin    = "admin"
)

// User is a registered account (email/password and/or Google OAuth).
type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Email        string             `bson:"email" json:"email"`
	PasswordHash string             `bson:"passwordHash,omitempty" json:"-"`
	GoogleID     string             `bson:"googleId,omitempty" json:"-"`
	Role         string             `bson:"role" json:"role"`
	Name         string             `bson:"name" json:"name"`
	CreatedAt    time.Time          `bson:"createdAt" json:"createdAt"`
}
