package controllers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"github.com/google/uuid"
	"yt_backend/db"
	"yt_backend/models"
)

func AddVideoToWatchHistory(c *gin.Context) {
	var input struct {
		VideoID string `json:"video_id"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	watchEntry := models.VideoWatchEntry{
		ID:        uuid.New().String(),
		UserID:    userID.(string),
		VideoID:   input.VideoID,
		WatchedAt: time.Now(),
	}

	watchHistoryCollection := db.GetCollection("video_watches")
	result, err := watchHistoryCollection.InsertOne(context.TODO(), watchEntry)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Video added to watch history", "id": result.InsertedID})
}

func GetWatchHistory(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var page int = 1
	if pageStr := c.Query("page"); pageStr != "" {
		page, _ = strconv.Atoi(pageStr)
	}

	var limit int = 20
	if limitStr := c.Query("limit"); limitStr != "" {
		limit, _ = strconv.Atoi(limitStr)
	}

	skip := (page - 1) * limit

	watchHistoryCollection := db.GetCollection("video_watches")
	
	pipeline := []bson.M{
			{
				"$match": bson.M{"user_id": userID},
			},
			{
				"$sort": bson.M{"watched_at": -1},
			},
			{
				"$skip": int64(skip),
			},
			{
				"$limit": int64(limit),
			},
		}

	cursor, err := watchHistoryCollection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var entries []models.VideoWatchEntry
	if err := cursor.All(context.TODO(), &entries); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"watch_history": entries, "page": page, "limit": limit})
}

func DeleteVideoFromWatchHistory(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	videoID := c.Param("video_id")

	watchHistoryCollection := db.GetCollection("video_watches")
	result, err := watchHistoryCollection.DeleteOne(context.TODO(),
		bson.M{"user_id": userID, "video_id": videoID},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Video not found in watch history"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Video removed from watch history"})
}
