package models

import (
	"slices"
	"time"
)

type Playlist struct {
	ID        string    `json:"id" bson:"_id"`
	UserID    string    `json:"userId" bson:"userId"`
	Name      string    `json:"name" bson:"name"`
	IsPublic  bool      `json:"isPublic" bson:"isPublic"`
	VideoIDs  []string  `json:"videoIds" bson:"videoIds"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

// AddVideo adds a video to the playlist
func (p *Playlist) AddVideo(videoID string) {
	if slices.Contains(p.VideoIDs, videoID) {
		return
	}
	p.VideoIDs = append(p.VideoIDs, videoID)
	p.UpdatedAt = time.Now()
}

// HasVideo checks if a video exists in the playlist
func (p *Playlist) HasVideo(videoID string) bool {
	return slices.Contains(p.VideoIDs, videoID)
}

// RemoveVideo removes a video from the playlist
func (p *Playlist) RemoveVideo(videoID string) {
	p.VideoIDs = slices.Delete(p.VideoIDs, slices.Index(p.VideoIDs, videoID), slices.Index(p.VideoIDs, videoID)+1)
	p.UpdatedAt = time.Now()
}

// GetVideoCount returns the number of videos in the playlist
func (p *Playlist) GetVideoCount() int {
	return len(p.VideoIDs)
}
