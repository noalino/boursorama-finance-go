package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/noalino/boursorama-finance-go/internal/api"
)

func main() {
	router := api.Router{Engine: gin.Default()}

	router.GET("/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "OK",
		})
	})

	router.RegisterHandlers()

	router.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
