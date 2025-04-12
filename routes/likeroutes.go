package routes

import (
	"yt_backend/controllers"
	"yt_backend/middleware"

	"github.com/gin-gonic/gin"
)

func LikeRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/video/:videoID/like", middleware.AuthMiddleware(), controllers.LikeVideo)
	incomingRoutes.DELETE("/video/:videoID/like", middleware.AuthMiddleware(), controllers.RemoveLike)
	incomingRoutes.GET("/video/:videoID/likes", controllers.CountVideoLikes)
}
