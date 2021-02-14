package api

import (
    "log"
    "net/http"
    "time"

    "github.com/gin-gonic/gin"

    "github.com/benoitgelineau/go-fetch-quotes/internal/utils"
)

func RegisterHandlers(router *gin.Engine) {
    router.GET("/search", func(c *gin.Context) {
        q := c.Query("q")

        results := utils.ScrapeSearchResult(q);
        c.JSON(http.StatusOK, results)
    })

    router.GET("/quotes/:symbol", func(c *gin.Context) {
        symbol := c.Param("symbol")
        // https://github.com/gin-gonic/gin#custom-validators
        now := time.Now()
        lastMonth := now.AddDate(0,-1,0)
        // Default start date = a month from now
        startDate := c.DefaultQuery("startDate", lastMonth.Format(services.LayoutISO))
        startDateAsTime, err := time.Parse(services.LayoutISO, startDate)
        if err != nil {
            log.Fatal(err)
        }
        // Default duration = 3 months
        duration := c.DefaultQuery("duration", "3M")
        // Default period = daily
        period := c.DefaultQuery("period", "1")

        quotes := utils.GetQuotes(symbol, startDateAsTime, duration, period)

        c.JSON(http.StatusOK, quotes)
    })
}
