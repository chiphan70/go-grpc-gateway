package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents a user document in MongoDB
type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string             `bson:"name" json:"name"`
	Email     string             `bson:"email" json:"email"`
	Phone     string             `bson:"phone" json:"phone"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

// NewUser creates a new user with timestamps
func NewUser(name, email, phone string) *User {
	now := time.Now()
	return &User{
		Name:      name,
		Email:     email,
		Phone:     phone,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
