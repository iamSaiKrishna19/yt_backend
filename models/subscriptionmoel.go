package models

import "time"

type Subscription struct {
	ID        string `json:"id" bson:"_id"`
	ChannelName Channel  `json:"channelName" bson:"channelName"`
	Subscribers User `json:"subscribers" bson:"subscribers"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

