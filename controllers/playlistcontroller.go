package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"go.mongodb.org/mongo-driver/bson"

	"yt_backend/db"
	"yt_backend/models"
)

func CreatePlaylist(c *gin.Context) {
	var input struct {
		Name     string `json:"name" binding:"required"`
		IsPublic *bool  `json:"isPublic"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Set default value for IsPublic if not provided
	isPublic := true
	if input.IsPublic != nil {
		isPublic = *input.IsPublic
	}

	// Create new playlist
	playlist := &models.Playlist{
		ID:        uuid.New().String(),
		UserID:    userID.(string),
		Name:      input.Name,
		IsPublic:  isPublic,
		VideoIDs:  []string{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Insert into database
	collection := db.GetCollection("playlists")
	_, err := collection.InsertOne(context.Background(), playlist)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create playlist"})
		return
	}

	c.JSON(http.StatusOK, playlist)
}

func AddToPlaylist(c *gin.Context) {
	playlistID := c.Param("playlistId")
	if playlistID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Playlist ID is required"})
		return
	}

	var input struct {
		VideoID string `json:"videoId" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get playlist from database
	collection := db.GetCollection("playlists")
	var playlist models.Playlist
	if err := collection.FindOne(context.Background(), bson.M{"_id": playlistID}).Decode(&playlist); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Playlist not found"})
		return
	}

	// Check if video already exists
	if playlist.HasVideo(input.VideoID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Video already exists in playlist"})
		return
	}

	// Update playlist in database using $push operator
	update := bson.M{
		"$push": bson.M{
			"videoIds": input.VideoID,
		},
		"$set": bson.M{
			"updatedAt": time.Now(),
		},
	}

	_, err := collection.UpdateOne(context.Background(), bson.M{"_id": playlistID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add video to playlist"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Video added to playlist successfully"})
}

func DeletePlaylist(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	playlistID := c.Param("playlistId")
	if playlistID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Playlist ID is required"})
		return
	}

	// Get playlist from database
	collection := db.GetCollection("playlists")
	var playlist models.Playlist
	if err := collection.FindOne(context.Background(), bson.M{"_id": playlistID}).Decode(&playlist); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Playlist not found"})
		return
	}

	// Check if user owns the playlist
	if playlist.UserID != userID.(string) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authorized to delete playlist"})
		return
	}

	// Delete playlist from database
	_, err := collection.DeleteOne(context.Background(), bson.M{"_id": playlistID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete playlist"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Playlist deleted successfully"})
}

func RemoveFromPlaylist(c *gin.Context) {
	playlistID := c.Param("playlistId")
	if playlistID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Playlist ID is required"})
		return
	}

	var input struct {
		VideoID string `json:"videoId" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get playlist from database
	collection := db.GetCollection("playlists")
	var playlist models.Playlist
	if err := collection.FindOne(context.Background(), bson.M{"_id": playlistID}).Decode(&playlist); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Playlist not found"})
		return
	}

	// Check if video exists
	if !playlist.HasVideo(input.VideoID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Video not found in playlist"})
		return
	}

	// Update playlist in database using $pull operator
	update := bson.M{
		"$pull": bson.M{
			"videoIds": input.VideoID,
		},
		"$set": bson.M{
			"updatedAt": time.Now(),
		},
	}

	_, err := collection.UpdateOne(context.Background(), bson.M{"_id": playlistID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove video from playlist"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Video removed from playlist successfully"})
}
