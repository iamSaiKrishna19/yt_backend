package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"yt_backend/db"
)

func main() {
	db.ConnectDB();
	router := gin.Default()
	router.GET("/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello World!",
		})
	})


	router.Run(":8080")// listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
