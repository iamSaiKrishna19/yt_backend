package models

import "time"

type VideoComment struct {
	ID        string    `json:"id" bson:"_id"`
	Content   string    `json:"content" bson:"content"`
	Owner     User      `json:"owner"`
	VComment  Video     `json:"vcomment" bson:"vcomment"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}
