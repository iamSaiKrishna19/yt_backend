package routes

import (
	"yt_backend/controllers"
	"yt_backend/middleware"

	"github.com/gin-gonic/gin"
)

func SetupWatchHistoryRoutes(incomingRoutes *gin.Engine) {
	watchHistory := incomingRoutes.Group("/watch-history")
	{
		watchHistory.POST("/add", middleware.AuthMiddleware(), controllers.AddVideoToWatchHistory)
		watchHistory.GET("/history", middleware.AuthMiddleware(), controllers.GetWatchHistory)
		watchHistory.DELETE("/delete/:video_id", middleware.AuthMiddleware(), controllers.DeleteVideoFromWatchHistory)
	}
}
