package routes

import (
	"yt_backend/controllers"
	"yt_backend/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/users/signup", controllers.SignUp)
	incomingRoutes.POST("/users/login", controllers.Login)
	incomingRoutes.PUT("/users/change-password", middleware.AuthMiddleware(), controllers.ChangePassword)
	incomingRoutes.PUT("/users/create-channel", middleware.AuthMiddleware(), controllers.CreateChannel)
	incomingRoutes.GET("/users/subscribed-to-channel", middleware.AuthMiddleware(), controllers.SubscribedToChannel)
}
