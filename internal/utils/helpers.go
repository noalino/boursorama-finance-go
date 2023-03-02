package utils

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const (
	BASE_URL  = "https://www.boursorama.com"
	LayoutISO = "02/01/2006"
)

var DefaultDurations = []string{"1M", "2M", "3M", "4M", "5M", "6M", "7M", "8M", "9M", "10M", "11M", "1Y", "2Y", "3Y"}
var DefaultPeriods = []string{"1", "7", "30", "365"}

func getQuotesUrl(symbol string, startDate time.Time, duration string, period string, page int) string {
	if page == 1 {
		return BASE_URL + "/_formulaire-periode/?symbol=" + strings.ToUpper(symbol) + "&historic_search[startDate]=" + startDate.Format(LayoutISO) + "&historic_search[duration]=" + duration + "&historic_search[period]=" + period
	} else {
		return BASE_URL + "/_formulaire-periode/page-" + strconv.Itoa(page) + "?symbol=" + strings.ToUpper(symbol) + "&historic_search[startDate]=" + startDate.Format(LayoutISO) + "&historic_search[duration]=" + duration + "&historic_search[period]=" + period
	}
}

func getSearchUrl(searchValue string) string {
	return BASE_URL + "/recherche/_instruments/" + searchValue
}

func ValidateInput(input string) string {
	return url.QueryEscape(strings.TrimSpace(input))
}

func getHTMLDocument(url string) (*goquery.Document, error) {
	// Request the HTML page
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func getMaxPages(view *goquery.Document) int {
	page, err := strconv.Atoi(view.Find("span.c-pagination__content").Last().Text())
	if err != nil {
		return 1
	}
	return page
}

func contains(values []string, query string) bool {
	for _, value := range values {
		if value == query {
			return true
		}
	}
	return false
}
