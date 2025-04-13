package models

import "time"

type VideoWatchEntry struct {
	ID        string    `json:"id" bson:"_id" validate:"required"`
	UserID    string    `json:"user_id" bson:"user_id" validate:"required"`
	VideoID   string    `json:"video_id" bson:"video_id" validate:"required"`
	Video     Video     `json:"video" bson:"video"`
	WatchedAt time.Time `json:"watched_at" bson:"watched_at"`
}
