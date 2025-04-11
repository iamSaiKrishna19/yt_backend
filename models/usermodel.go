package models

import "time"

type User struct {
	ID           string    `json:"id" bson:"_id" validate:"required"`
	Username     string    `json:"username" bson:"username" validate:"required,min=3,max=20"`
	ChannelName  Channel    `json:"channelName" bson:"channelName"`
	Email        string    `json:"email" bson:"email" validate:"required,email" lowercase:"true"`
	Password     string    `json:"password" bson:"password" validate:"required,min=8"`
	Avatar       string    `json:"avatar" bson:"avatar" validate:"omitempty,url"`
	CoverImage   string    `json:"coverImage" bson:"coverImage" validate:"omitempty,url"`
	RefreshToken string    `json:"refreshToken" bson:"refreshToken" validate:"omitempty"`
	CreatedAt    time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt" bson:"updatedAt"`
}
