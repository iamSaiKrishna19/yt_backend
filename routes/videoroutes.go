package routes

import (
	"yt_backend/controllers"
	"yt_backend/middleware"
	"github.com/gin-gonic/gin"
)

func VideoRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/videos/upload", middleware.AuthMiddleware(), controllers.UploadVideo)
	incomingRoutes.DELETE("/videos/:videoId", middleware.AuthMiddleware(), controllers.DeleteVideo)
	incomingRoutes.GET("/videos/:videoId/views", controllers.IncrementViews)
}
