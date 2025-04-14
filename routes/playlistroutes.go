package routes

import (
	"yt_backend/controllers"

	"github.com/gin-gonic/gin"
)

func PlaylistRoutes(incomingroutes *gin.Engine) {
	incomingroutes.POST("/create", controllers.CreatePlaylist)
	incomingroutes.POST("/add/:playlistId", controllers.AddToPlaylist)
	incomingroutes.POST("/remove/:playlistId", controllers.RemoveFromPlaylist)
	incomingroutes.DELETE("/:playlistId", controllers.DeletePlaylist)
}
