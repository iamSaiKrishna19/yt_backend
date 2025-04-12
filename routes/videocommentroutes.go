package routes

import (
	"yt_backend/controllers"
	"yt_backend/middleware"

	"github.com/gin-gonic/gin"
)

func VideoCommentRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/:videoId/post-comment", middleware.AuthMiddleware(), controllers.PostComment)
	incomingRoutes.DELETE("/:videoId/:commentId", middleware.AuthMiddleware(), controllers.DeleteComment)
	incomingRoutes.PUT("/video/:commentId",middleware.AuthMiddleware(),controllers.EditComment)
}
