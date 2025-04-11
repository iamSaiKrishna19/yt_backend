package controllers

import (
	"context"
	"net/http"
	"time"
	"yt_backend/db"
	"yt_backend/models"
	"yt_backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/mongo"
)

func UploadVideo(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Parse multipart form
	err := c.Request.ParseMultipartForm(100 << 20) // 100 MB max
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form"})
		return
	}

	// Get video file
	videoFile, err := c.FormFile("video")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Video file is required"})
		return
	}

	// Get title from form
	title := c.PostForm("title")
	if title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title is required"})
		return
	}

	// Get user details
	userCollection := db.GetCollection("users")
	var user models.User
	err = userCollection.FindOne(context.TODO(), bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Upload video to Cloudinary
	videoURL, err := utils.HandleVideoUpload(c.Request.Context(), videoFile, "videos")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload video to Cloudinary"})
		return
	}

	// Get video duration
	duration, err := utils.GetVideoDuration(videoFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get video duration"})
		return
	}

	// Create video document
	video := models.Video{
		ID:          uuid.New().String(),
		Title:       title,
		URL:         videoURL,
		Owner:       user,
		ChannelName: user.ChannelName,
		Duration:    duration,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Save video to database
	videoCollection := db.GetCollection("videos")
	_, err = videoCollection.InsertOne(context.TODO(), video)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save video"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Video uploaded successfully",
		"video":   video,
	})
}

func DeleteVideo(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	videoID := c.Param("videoId")
	if videoID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Video ID is required"})
		return
	}

	// Get video collection
	videoCollection := db.GetCollection("videos")

	// Find the video and verify ownership
	var video models.Video
	err := videoCollection.FindOne(context.TODO(), bson.M{"_id": videoID}).Decode(&video)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Video not found"})
		return
	}

	// Check if the user is the owner of the video
	if video.Owner.ID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to delete this video"})
		return
	}

	// Delete the video from Cloudinary
	err = utils.DeleteFromCloudinary(c.Request.Context(), video.URL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete video from Cloudinary"})
		return
	}

	// Delete the video from database
	_, err = videoCollection.DeleteOne(context.TODO(), bson.M{"_id": videoID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete video"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Video deleted successfully",
	})
}

func IncrementViews(c *gin.Context) {
	videoID := c.Param("videoId")
	if videoID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Video ID is required"})
		return
	}

	videoCollection := db.GetCollection("videos")

	var video models.Video
	err := videoCollection.FindOne(context.TODO(), bson.M{"_id": videoID}).Decode(&video)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Video not found"})
		return
	}

	_, err = videoCollection.UpdateOne(context.TODO(), bson.M{"_id": videoID}, bson.M{"$inc": bson.M{"views": 1}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to increment views"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Views incremented successfully",
	})
}
