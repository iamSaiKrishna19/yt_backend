package controllers

import (
	"context"
	"net/http"
	"time"
	"yt_backend/db"
	"yt_backend/models"

	// "yt_backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/mongo"
)

func PostComment(c *gin.Context) {
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

	content := c.PostForm("content")
	if content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Comment cannot be empty"})
	}

	// Get user details
	userCollection := db.GetCollection("users")
	var user models.User
	err := userCollection.FindOne(context.TODO(), bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	videoCollection := db.GetCollection("videos")
	var video models.Video
	err = videoCollection.FindOne(context.TODO(), bson.M{"_id": videoID}).Decode(&video)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Video not found"})
		return
	}

	comment := models.VideoComment{
		ID:        uuid.New().String(),
		Content:   content,
		Owner:     user,
		VComment:  video,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	commentCollection := db.GetCollection("videocomments")
	_, err = commentCollection.InsertOne(context.TODO(), comment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload comment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Comment uploaded successfully",
		"comment": comment,
	})
}

func DeleteComment(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	commentID := c.Param("commentID")
	if commentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Video ID is required"})
		return
	}

	commentCollection := db.GetCollection("videocomments")

	var comment models.VideoComment
	err := commentCollection.FindOne(context.TODO(), bson.M{"_id": commentID}).Decode(&comment)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}

	if comment.Owner.ID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to delete this comment"})
		return
	}

	_, err = commentCollection.DeleteOne(context.TODO(), bson.M{"_id": commentID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete comment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Comment deleted successfully",
	})

}

func EditComment(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	commentID := c.Param("commentID")
	if commentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Video ID is required"})
		return
	}

	content := c.PostForm("content")
	if content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Comment cannot be empty"})
	}

	commentCollection := db.GetCollection("videocomments")

	var comment models.VideoComment
	err := commentCollection.FindOne(context.TODO(), bson.M{"_id": commentID}).Decode(&comment)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}

	if comment.Owner.ID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to delete this comment"})
		return
	}

	update := bson.M{
		"$set": bson.M{
			"contect":   content,
			"updatedAt": time.Now(),
		},
	}

	_, err = commentCollection.UpdateOne(
		context.TODO(),
		bson.M{"_id": commentID},
		update,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update comment"})
		return
	}

	var editedComment models.User
	err = commentCollection.FindOne(context.TODO(), bson.M{"_id": commentID}).Decode(&editedComment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch updated user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Comment edited successfully",
		"comment": editedComment,
	})

}
