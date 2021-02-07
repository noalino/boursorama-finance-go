package services

import (
    "fmt"
    "log"
    "net/http"
    "strings"

    "github.com/PuerkitoBio/goquery"
)

type Asset struct {
    Code        string  `json:"code"`
    Name        string  `json:"name"`
    Category    string  `json:"category"`
    LastPrice   string  `json:"last_price"`
}

var searchResults []Asset

func ScrapeSearchResult(query string) []Asset {
    // Request the HTML page.
	res, err := http.Get("https://www.boursorama.com/recherche/ajax?query=" + query)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

    // Find the search results
	doc.Find(".search__list").First().Find(".search__list-link").Each(func(i int, s *goquery.Selection) {
        asset := Asset{}
		asset.Name = s.Find(".search__item-title").Text()
        link, ok := s.Attr("href")
        if !ok {
            log.Fatal("Unable to find link reference")
        }
        splittedLink := strings.Split(link, "/")
        asset.Code = splittedLink[len(splittedLink)-2]
        asset.Category = strings.Trim(s.Find(".search__item-content").Text(), " \n")
        asset.LastPrice = s.Find(".search__item-instrument .last").Text()
        searchResults = append(searchResults, asset)
	})

    return searchResults
}
