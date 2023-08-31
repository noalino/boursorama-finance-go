package utils

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	BASE_URL = "https://www.boursorama.com"
)

func getSearchUrl(searchValue string) string {
	return fmt.Sprintf("%s/recherche/_instruments/%s", BASE_URL, searchValue)
}

func getQuotesUrl(symbol string, startDate string, duration string, period string, page int) string {
	if page == 1 {
		return fmt.Sprintf("%s/_formulaire-periode/?symbol=%s&historic_search[startDate]=%s&historic_search[duration]=%s&historic_search[period]=%s", BASE_URL, strings.ToUpper(symbol), startDate, duration, period)
	} else {
		return fmt.Sprintf("%s/_formulaire-periode/page-%s?symbol=%s&historic_search[startDate]=%s&historic_search[duration]=%s&historic_search[period]=%s", BASE_URL, strconv.Itoa(page), strings.ToUpper(symbol), startDate, duration, period)
	}
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
