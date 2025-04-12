package controllers

import (
	"context"
	"net/http"
	"time"
	"yt_backend/db"
	"yt_backend/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// i have to do subscriber already exist or not
func Subscribe(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userCollection := db.GetCollection("users")
	var user models.User
	err := userCollection.FindOne(context.TODO(), bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	videoID := c.Param("videoId")
	if videoID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Video ID is required"})
		return
	}

	videoCollection := db.GetCollection("videos")

	var video models.Video
	err = videoCollection.FindOne(context.TODO(), bson.M{"_id": videoID}).Decode(&video)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Video not found"})
		return
	}

	subscriptionCollection := db.GetCollection("subscriptions")

	// Check if subscription already exists
	filter := bson.M{
		"subscribers._id": userID,
		"channelName._id": video.ChannelName.ID,
	}

	var existingSubscription models.Subscription
	err = subscriptionCollection.FindOne(context.TODO(), filter).Decode(&existingSubscription)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You are already subscribed to this channel"})
		return
	} else if err != mongo.ErrNoDocuments {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking existing subscription"})
		return
	}

	subscription := models.Subscription{
		ID:          uuid.New().String(),
		ChannelName: video.ChannelName,
		Subscribers: user,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	_, err = subscriptionCollection.InsertOne(context.TODO(), subscription)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save subscription"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Subscribed successfully"})
}

func Unsubscribe(c *gin.Context) {
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

	videoCollection := db.GetCollection("videos")

	var video models.Video
	err := videoCollection.FindOne(context.TODO(), bson.M{"_id": videoID}).Decode(&video)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Video not found"})
		return
	}

	subscriptionCollection := db.GetCollection("subscriptions")

	_, err = subscriptionCollection.DeleteOne(context.TODO(), bson.M{"subscribers._id": userID, "channelName._id": video.ChannelName.ID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unsubscribe"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Unsubscribed successfully"})
}

func CountSubscribers(c *gin.Context) {
	channelID := c.Param("channelId")
	if channelID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Channel ID is required"})
		return
	}

	subscriptionCollection := db.GetCollection("subscriptions")

	// Create aggregation pipeline
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"channelName._id": channelID,
			},
		},
		{
			"$group": bson.M{
				"_id": "$channelName._id",
				"subscriberCount": bson.M{
					"$sum": 1,
				},
			},
		},
	}

	// Execute aggregation
	cursor, err := subscriptionCollection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count subscribers"})
		return
	}
	defer cursor.Close(context.TODO())

	// Get the result
	var result []bson.M
	if err = cursor.All(context.TODO(), &result); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process subscriber count"})
		return
	}

	// If no subscribers found, return 0
	count := 0
	if len(result) > 0 {
		count = int(result[0]["subscriberCount"].(int32))
	}

	c.JSON(http.StatusOK, gin.H{
		"channelId":       channelID,
		"subscriberCount": count,
	})
}

// DemoPipelineOperations demonstrates different MongoDB pipeline operations
// func DemoPipelineOperations(c *gin.Context) {
// 	subscriptionCollection := db.GetCollection("subscriptions")

// 	// Example 1: Basic Count with Grouping
// 	basicPipeline := []bson.M{
// 		{
// 			"$group": bson.M{
// 				"_id": "$channelName._id", // Group by channel ID
// 				"totalSubscribers": bson.M{
// 					"$sum": 1, // Count documents in each group
// 				},
// 			},
// 		},
// 	}

// 	// Example 2: Count with Date Filtering
// 	dateFilterPipeline := []bson.M{
// 		{
// 			"$match": bson.M{ // Filter documents first
// 				"createdAt": bson.M{
// 					"$gte": time.Now().AddDate(0, -1, 0), // Last 1 month
// 				},
// 			},
// 		},
// 		{
// 			"$group": bson.M{
// 				"_id": "$channelName._id",
// 				"recentSubscribers": bson.M{
// 					"$sum": 1,
// 				},
// 			},
// 		},
// 	}

// 	// Example 3: Complex Aggregation with Multiple Stages
// 	complexPipeline := []bson.M{
// 		{
// 			"$match": bson.M{ // Stage 1: Filter
// 				"createdAt": bson.M{
// 					"$exists": true,
// 				},
// 			},
// 		},
// 		{
// 			"$group": bson.M{ // Stage 2: Group
// 				"_id": "$channelName._id",
// 				"totalSubscribers": bson.M{
// 					"$sum": 1,
// 				},
// 				"firstSubscription": bson.M{
// 					"$min": "$createdAt", // Get earliest subscription date
// 				},
// 				"lastSubscription": bson.M{
// 					"$max": "$createdAt", // Get latest subscription date
// 				},
// 			},
// 		},
// 		{
// 			"$project": bson.M{ // Stage 3: Reshape the output
// 				"channelId":        "$_id",
// 				"totalSubscribers": 1,
// 				"subscriptionDuration": bson.M{
// 					"$subtract": []interface{}{
// 						"$lastSubscription",
// 						"$firstSubscription",
// 					},
// 				},
// 			},
// 		},
// 		{
// 			"$sort": bson.M{ // Stage 4: Sort
// 				"totalSubscribers": -1, // Descending order
// 			},
// 		},
// 		{
// 			"$limit": 10, // Stage 5: Limit results
// 		},
// 	}

// 	// Execute the complex pipeline
// 	cursor, err := subscriptionCollection.Aggregate(context.TODO(), complexPipeline)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute pipeline"})
// 		return
// 	}
// 	defer cursor.Close(context.TODO())

// 	// Get the results
// 	var results []bson.M
// 	if err = cursor.All(context.TODO(), &results); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process results"})
// 		return
// 	}

// 	// Common Pipeline Stages Explained:
// 	// 1. $match: Filters documents
// 	//    Example: {"$match": {"field": "value"}}

// 	// 2. $group: Groups documents by a field
// 	//    Example: {"$group": {"_id": "$field", "count": {"$sum": 1}}}

// 	// 3. $project: Reshapes documents
// 	//    Example: {"$project": {"newField": "$oldField"}}

// 	// 4. $sort: Sorts documents
// 	//    Example: {"$sort": {"field": 1}} // 1 for ascending, -1 for descending

// 	// 5. $limit: Limits number of documents
// 	//    Example: {"$limit": 10}

// 	// 6. $skip: Skips documents
// 	//    Example: {"$skip": 5}

// 	// 7. $unwind: Deconstructs an array field
// 	//    Example: {"$unwind": "$arrayField"}

// 	// Common Aggregation Operators:
// 	// $sum: Adds values
// 	// $avg: Calculates average
// 	// $min: Finds minimum value
// 	// $max: Finds maximum value
// 	// $first: Gets first value in group
// 	// $last: Gets last value in group
// 	// $push: Creates array of values
// 	// $addToSet: Creates array of unique values

// 	c.JSON(http.StatusOK, gin.H{
// 		"message": "Pipeline demonstration",
// 		"results": results,
// 	})
// }
