package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/benoitgelineau/boursorama-finance-go/internal/api"
)

func main() {
	router := gin.Default()

	router.GET("/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "OK",
		})
	})

	api.RegisterHandlers(router)

	router.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
