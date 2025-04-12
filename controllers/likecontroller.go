package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"
	"yt_backend/db"
	"yt_backend/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func LikeVideo(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	videoID := c.Param("videoID")
	if videoID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Video ID is required"})
		return
	}

	likeCollection := db.GetCollection("likes")

	// Check if like already exists
	filter := bson.M{
		"owner._id": userID,
		"vlike._id": videoID,
	}

	userCollection := db.GetCollection("users")
	var user models.User
	err := userCollection.FindOne(context.TODO(), bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	videoCollection := db.GetCollection("videos")
	var video models.Video
	err = videoCollection.FindOne(context.TODO(), bson.M{"_id": videoID}).Decode(&video)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Video not found"})
		return
	}

	var existingLike models.Like
	err = likeCollection.FindOne(context.TODO(), filter).Decode(&existingLike)

	fmt.Println(err)

	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You already liked this video"})
		return
	} else if err != mongo.ErrNoDocuments {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking existing like"})
		return
	}

	// Like does not exist, create one
	newLike := models.Like{
		ID:        uuid.New().String(),
		Owner:     user,
		VLike:     video,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err = likeCollection.InsertOne(context.TODO(), newLike)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to like video"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Video liked successfully"})
}

func RemoveLike(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	videoID := c.Param("videoID")
	if videoID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Video ID is required"})
		return
	}

	likeCollection := db.GetCollection("likes")

	// Check if like exists before attempting to remove
	filter := bson.M{
		"owner._id": userID,
		"vlike._id": videoID,
	}

	var existingLike models.Like
	err := likeCollection.FindOne(context.TODO(), filter).Decode(&existingLike)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "You haven't liked this video"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking like status"})
		}
		return
	}

	// Remove the like
	result, err := likeCollection.DeleteOne(context.TODO(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove like"})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Like not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Like removed successfully"})
}

func CountVideoLikes(c *gin.Context) {
	videoID := c.Param("videoID")
	if videoID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Video ID is required"})
		return
	}

	likeCollection := db.GetCollection("likes")

	// Create aggregation pipeline to count likes
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"vlike._id": videoID,
			},
		},
		{
			"$group": bson.M{
				"_id": "$vlike._id",
				"likeCount": bson.M{
					"$sum": 1,
				},
			},
		},
	}

	// Execute aggregation
	cursor, err := likeCollection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count likes"})
		return
	}
	defer cursor.Close(context.TODO())

	// Get the result
	var result []bson.M
	if err = cursor.All(context.TODO(), &result); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process like count"})
		return
	}

	// If no likes found, return 0
	count := 0
	if len(result) > 0 {
		count = int(result[0]["likeCount"].(int32))
	}

	c.JSON(http.StatusOK, gin.H{
		"videoId":   videoID,
		"likeCount": count,
	})
}
