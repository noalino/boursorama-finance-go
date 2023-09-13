package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/noalino/boursorama-finance-go/internal/lib"
	options "github.com/noalino/boursorama-finance-go/internal/lib/options/get"
)

type Router struct {
	*gin.Engine
}

func (router *Router) RegisterHandlers() {
	router.GET("/search", func(c *gin.Context) {
		page, err := strconv.ParseUint(c.DefaultQuery("page", "1"), 10, 16)
		if err != nil {
			handleBadRequest(c, err)
			return
		}

		query := lib.SearchQuery{
			Value: c.Query("q"),
			Page:  uint16(page),
		}

		results, err := lib.Search(query)
		if err != nil {
			handleBadRequest(c, err)
			return
		}

		c.JSON(http.StatusOK, results)
	})

	router.GET("/quotes/:symbol", func(c *gin.Context) {
		query := lib.GetQuery{
			Symbol:   c.Param("symbol"),
			From:     c.DefaultQuery("startDate", options.DefaultFrom().String()),
			Duration: c.DefaultQuery("duration", options.DefaultDuration.String()),
			Period:   c.DefaultQuery("period", options.DefaultPeriod.String()),
		}

		quotes, err := lib.Get(query)
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
