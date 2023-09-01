package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/noalino/boursorama-finance-go/internal/options"
	"github.com/noalino/boursorama-finance-go/internal/utils"
)

type Router struct {
	*gin.Engine
}

func (router *Router) RegisterHandlers() {
	router.GET("/search", func(c *gin.Context) {
		query := utils.SearchQuery{Value: c.Query("q")}

		results, err := utils.ScrapeSearchResult(query)
		if err != nil {
			handleBadRequest(c, err)
			return
		}

		c.JSON(http.StatusOK, results)
	})

	router.GET("/quotes/:symbol", func(c *gin.Context) {
		query := utils.QuotesQuery{
			Symbol: c.Param("symbol"),
			From:     c.DefaultQuery("startDate", options.DefaultFrom().String()),
			Duration: c.DefaultQuery("duration", options.DefaultDuration.String()),
			Period:   c.DefaultQuery("period", options.DefaultPeriod.String()),
		}

		quotes, err := utils.GetQuotes(query)
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
