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

func SanitizeUrlInput(input string) string {
	return url.QueryEscape(strings.TrimSpace(input))
}

func GetHTMLDocument(url string) (*goquery.Document, error) {
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

func GetMaxPages(view *goquery.Document) uint16 {
	page, err := strconv.ParseUint(view.Find("span.c-pagination__content").Last().Text(), 10, 16)
	if err != nil {
		return 1
	}
	return uint16(page)
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
