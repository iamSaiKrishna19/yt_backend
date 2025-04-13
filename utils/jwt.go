package utils

import (
	"context"
	"os"
	"time"
	"yt_backend/db"
	"yt_backend/models"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateToken(userID string) (string, error) {
	for {
		// Create the claims
		claims := Claims{
			userID,
			jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // Token expires in 24 hours
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				NotBefore: jwt.NewNumericDate(time.Now()),
			},
		}

		// Create the token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		// Sign the token with the secret key
		secretKey := os.Getenv("JWT_SECRET")
		if secretKey == "" {
			secretKey = "your-secret-key" // Default key, should be changed in production
		}

		signedToken, err := token.SignedString([]byte(secretKey))
		if err != nil {
			return "", err
		}

		// Check if token is blacklisted
		blacklistCollection := db.GetCollection("token_blacklist")
		var blacklistEntry models.TokenBlacklist
		filter := bson.M{"token": signedToken}
		err = blacklistCollection.FindOne(context.TODO(), filter).Decode(&blacklistEntry)
		if err == mongo.ErrNoDocuments {
			// Token is not blacklisted, we can use it
			return signedToken, nil
		} else if err != nil {
			// Handle other errors
			return "", err
		}

		// Token is blacklisted, generate a new one
		continue
	}
}

func VerifyToken(tokenString string) (*Claims, error) {
	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		secretKey := os.Getenv("JWT_SECRET")
		if secretKey == "" {
			secretKey = "your-secret-key" // Default key, should be changed in production
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	// Validate the token and return the claims
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}
