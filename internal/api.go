package api

import (
    "log"
    "net/http"
    "github.com/gin-gonic/gin"
)

func Init() {
	router := gin.Default()
	router.GET("/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "OK",
		})
	})

    router.GET("/search", func(c *gin.Context) {
        q := c.Query("q")
        log.Printf("/search query: %s", q)
        // TODO Check empty query, return 404? 204 No content?
        c.JSON(http.StatusOK, gin.H{
            "data": "[]",
        })
    })

    router.GET("/quotes/:code", func(c *gin.Context) {
        code := c.Param("code")
        // US date format
        // TODO choose default based on interval & date of request
        from := c.DefaultQuery("from", "2020-06-02")
        // TODO Default to date of request
        to := c.DefaultQuery("to", "2020-06-02")
        // TODO Default to max interval possible (based on dates), or 1M default
        interval := c.DefaultQuery("interval", "1M")
        log.Printf("/quotes code: %s, from: %s, to: %s, interval: %s", code, from, to, interval)
        c.JSON(http.StatusOK, gin.H{
            "data": "[]",
        })
    })

	router.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
