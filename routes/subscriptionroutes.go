package routes

import (
	"yt_backend/controllers"
	"yt_backend/middleware"

	"github.com/gin-gonic/gin"
)

func SubscriptionRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/subscribe/:videoId", middleware.AuthMiddleware(), controllers.Subscribe)
	incomingRoutes.DELETE("/unsubscribe/:videoId", middleware.AuthMiddleware(), controllers.Unsubscribe)
	incomingRoutes.GET("/subscribers/count/:channelId", controllers.CountSubscribers)
}
