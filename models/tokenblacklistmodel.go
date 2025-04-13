package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// TokenBlacklist represents a blacklisted token
type TokenBlacklist struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Token     string            `bson:"token" json:"token"`
	CreatedAt time.Time         `bson:"created_at" json:"created_at"`
}
