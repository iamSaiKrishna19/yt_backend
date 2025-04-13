package routes

import (
	"yt_backend/controllers"
	"yt_backend/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	userRoutes := incomingRoutes.Group("/users")
	{
		userRoutes.POST("/signup", controllers.SignUp)
		userRoutes.POST("/login", controllers.Login)
		userRoutes.POST("/logout", middleware.AuthMiddleware(), controllers.Logout)
		userRoutes.PUT("/change-password", middleware.AuthMiddleware(), controllers.ChangePassword)
		userRoutes.PUT("/create-channel", middleware.AuthMiddleware(), controllers.CreateChannel)
		userRoutes.GET("/subscribed-to-channel", middleware.AuthMiddleware(), controllers.SubscribedToChannel)
	}
}
