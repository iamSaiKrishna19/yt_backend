package models

import "time"

type Channel struct {
	ID        string `json:"id" bson:"_id"`
	ChannelName string `json:"channelName" bson:"channelName"`
	CreatedAt   time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt" bson:"updatedAt"`
}
