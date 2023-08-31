package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/noalino/boursorama-finance-go/internal/options"
	"github.com/noalino/boursorama-finance-go/internal/utils"
)

func RegisterHandlers(router *gin.Engine) {
	router.GET("/search", func(c *gin.Context) {
		q := c.Query("q")
		if q == "" {
			handleBadRequest(c, "Missing query value")
			return
		}

		results, err := utils.ScrapeSearchResult(q)
		if err != nil {
			handleBadRequest(c, err)
			return
		}
		c.JSON(http.StatusOK, results)
	})

	router.GET("/quotes/:symbol", func(c *gin.Context) {
		symbol := c.Param("symbol")
		startDate := c.DefaultQuery("startDate", options.DefaultFrom().String())
		duration := c.DefaultQuery("duration", options.DefaultDuration.String())
		period := c.DefaultQuery("period", options.DefaultPeriod.String())

		quotes, err := utils.GetQuotes(symbol, startDate, duration, period)
		if err != nil {
			handleBadRequest(c, err)
			return
		}

		c.JSON(http.StatusOK, quotes)
	})
}

func handleBadRequest(c *gin.Context, message interface{}) {
	c.JSON(http.StatusBadRequest, gin.H{
		"status":  http.StatusBadRequest,
		"message": message,
	})
}
