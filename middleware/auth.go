package middleware

import (
	"context"
	"net/http"
	"strings"

	"yt_backend/db"
	"yt_backend/models"
	"yt_backend/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// Check if the header starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		// Extract the token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Verify the token
		claims, err := utils.VerifyToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Check if token is blacklisted ////added only these
		blacklistCollection := db.GetCollection("token_blacklist")
		var blacklistEntry models.TokenBlacklist
		filter := bson.M{"token": tokenString}
		err = blacklistCollection.FindOne(context.TODO(), filter).Decode(&blacklistEntry)
		if err == nil {
			// Token found in blacklist
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token has been revoked"})
			c.Abort()
			return
		} else if err != mongo.ErrNoDocuments {
			// Handle other errors
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			c.Abort()
			return
		}
		// till these lines the token is verified and the user is authenticated

		// Add the user ID to the context
		c.Set("user_id", claims.UserID)
		c.Next()
	}
}
