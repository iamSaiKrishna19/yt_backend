package models

import "time"

type Video struct {
	ID           string    `json:"id" bson:"_id" validate:"required"`
	Title       string    `json:"title"`
	URL         string    `json:"url"`
	Owner       User      `json:"owner"`
	ChannelName Channel   `json:"channel_name"`
	Views       int       `json:"views" default:"0"`
	Duration    string    `json:"duration"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
