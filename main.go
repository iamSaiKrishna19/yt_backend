package main

import (
	"net/http"

	"yt_backend/db"
	"yt_backend/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	db.ConnectDB()
	router := gin.Default()
	router.GET("/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello World!",
		})
	})
	routes.UserRoutes(router)
	routes.VideoRoutes(router)
	routes.SubscriptionRoutes(router)
	routes.VideoCommentRoutes(router)
	routes.LikeRoutes(router)

	router.Run(":8080") // "localhost:8080"
}
