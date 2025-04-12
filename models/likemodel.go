package models

import "time"

type Like struct {
	ID        string    `json:"id" bson:"_id"`
	Owner     User      `json:"owner"`
	VLike     Video     `json:"vlike" bson:"vlike"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}
