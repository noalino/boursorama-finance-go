package utils

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const (
	BASE_URL  = "https://www.boursorama.com"
	LayoutISO = "02/01/2006"
)

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

func PrintlnOrVoid(condition bool, args ...any) {
	if condition {
		fmt.Println(args...)
	}
}

func PrintfOrVoid(condition bool, text string, args ...any) {
	if condition {
		fmt.Printf(text, args...)
	}
}

func IsDataFromPipe() bool {
	fi, _ := os.Stdin.Stat()
	return (fi.Mode() & os.ModeCharDevice) == 0
}
