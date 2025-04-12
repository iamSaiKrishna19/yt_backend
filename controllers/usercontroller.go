package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"
	"yt_backend/db"
	"yt_backend/models"
	"yt_backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func HashPass(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func SignUp(c *gin.Context) {
	// Parse multipart form
	err := c.Request.ParseMultipartForm(10 << 20) // 10 MB max
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form"})
		return
	}

	// Create user struct
	user := models.User{
		ID:       uuid.New().String(), // Generate unique ID
		Username: c.PostForm("username"),
		Email:    c.PostForm("email"),
		Password: c.PostForm("password"),
	}

	// Validate required fields
	if user.Username == "" || user.Email == "" || user.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username, email, and password are required"})
		return
	}

	// Check if user already exists
	collection := db.GetCollection("users")
	filter := bson.M{
		"$or": []bson.M{
			{"Username": user.Username},
			{"email": user.Email},
		},
	}

	err = collection.FindOne(context.TODO(), filter).Err()
	fmt.Println(err)

	if err == nil {
		// Found a user – duplicate
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username or Email already exists"})
		return
	}

	if err != mongo.ErrNoDocuments {
		// Real database error
		fmt.Println("Database error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}

	// Hash password
	hashedPassword, err := HashPass(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	user.Password = hashedPassword

	// Set timestamps
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	// Handle avatar upload (optional)
	avatarFile, err := c.FormFile("avatar")
	if err == nil && avatarFile != nil {
		avatarURL, err := utils.HandleImageUpload(c.Request.Context(), avatarFile, "avatar")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload avatar"})
			return
		}
		user.Avatar = avatarURL
	} else {
		user.Avatar = ""
	}

	// Handle cover image upload (optional)
	coverFile, err := c.FormFile("coverImage")
	if err == nil && coverFile != nil {
		coverURL, err := utils.HandleImageUpload(c.Request.Context(), coverFile, "cover")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload cover image"})
			return
		}
		user.CoverImage = coverURL
	} else {
		user.CoverImage = ""
	}

	// Save user to MongoDB
	_, err = collection.InsertOne(context.Background(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save user to database"})
		return
	}

	// Return success response
	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"user":    user,
	})
}

func Login(c *gin.Context) {
	var loginData struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.BindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// Validate that either username or email is provided
	if loginData.Username == "" && loginData.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username or email is required"})
		return
	}

	// Create filter based on provided credentials
	var filter bson.M
	if loginData.Username != "" {
		filter = bson.M{"username": loginData.Username}
	} else {
		filter = bson.M{"email": loginData.Email}
	}

	// Find user in database
	collection := db.GetCollection("users")
	var user models.User
	err := collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username/email or password"})
		return
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
		"user":    user,
	})
}

func ChangePassword(c *gin.Context) {

	var changePass struct {
		CurrentPassword string `json:"currentPassword"`
		NewPassword     string `json:"newPassword"`
	}

	if err := c.BindJSON(&changePass); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// Validate passwords
	if changePass.CurrentPassword == "" || changePass.NewPassword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Current password and new password are required"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Find user in database
	collection := db.GetCollection("users")
	var user models.User
	err := collection.FindOne(context.TODO(), bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Verify current password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(changePass.CurrentPassword))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Current password is incorrect"})
		return
	}

	// Hash new password
	hashedPassword, err := HashPass(changePass.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash new password"})
		return
	}

	// Update password in database
	update := bson.M{
		"$set": bson.M{
			"password":  hashedPassword,
			"updatedAt": time.Now(),
		},
	}

	_, err = collection.UpdateOne(context.TODO(), bson.M{"_id": userID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Password updated successfully",
	})
}

func CreateChannel(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get channel name from form data
	channelName := c.PostForm("channelName")
	if channelName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Channel name is required"})
		return
	}

	// Create a new Channel instance
	channel := models.Channel{
		ID:          uuid.New().String(),
		ChannelName: channelName,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Insert channel into channels collection
	channelCollection := db.GetCollection("channels")
	_, err := channelCollection.InsertOne(context.TODO(), channel)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create channel"})
		return
	}

	// Update the user's channel reference
	userCollection := db.GetCollection("users")
	update := bson.M{
		"$set": bson.M{
			"channelName": channel,
			"updatedAt":   time.Now(),
		},
	}

	// Update the user document
	_, err = userCollection.UpdateOne(
		context.TODO(),
		bson.M{"_id": userID},
		update,
	)
	if err != nil {
		// If user update fails, we should clean up the channel
		_, _ = channelCollection.DeleteOne(context.TODO(), bson.M{"_id": channel.ID})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user channel"})
		return
	}

	// Get the updated user
	var updatedUser models.User
	err = userCollection.FindOne(context.TODO(), bson.M{"_id": userID}).Decode(&updatedUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch updated user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Channel created successfully",
		"user":    updatedUser,
	})
}

func SubscribedToChannel(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	subscriptionCollection := db.GetCollection("subscriptions")

	// Create pipeline to get subscribed channels
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"subscribers._id": userID,
			},
		},
		{
			"$lookup": bson.M{ // Join with channels collection
				"from":         "channels",
				"localField":   "channelName._id",
				"foreignField": "_id",
				"as":           "channelDetails",
			},
		},
		{
			"$unwind": "$channelDetails", // Flatten the channelDetails array
		},
		{
			"$project": bson.M{ // Shape the output
				"channelId":    "$channelDetails._id",
				"channelName":  "$channelDetails.channelName",
				"createdAt":    "$channelDetails.createdAt",
				"subscribedAt": "$createdAt", // When user subscribed
				"subscriberCount": bson.M{
					"$size": "$channelDetails.subscribers", // Count of subscribers
				},
			},
		},
		{
			"$sort": bson.M{ // Sort by subscription date
				"subscribedAt": -1,
			},
		},
	}

	cursor, err := subscriptionCollection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch subscribed channels"})
		return
	}
	defer cursor.Close(context.TODO())

	var subscribedChannels []bson.M
	if err := cursor.All(context.TODO(), &subscribedChannels); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process subscribed channels"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Subscribed channels retrieved successfully",
		"channels": subscribedChannels,
		"count":    len(subscribedChannels),
	})
}
