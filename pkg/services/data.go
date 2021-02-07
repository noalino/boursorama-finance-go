package services

import (
    "log"
    "net/http"
    "time"
    "strings"

    "github.com/PuerkitoBio/goquery"
)

type Asset struct {
    Symbol      string  `json:"symbol"`
    Name        string  `json:"name"`
    Category    string  `json:"category"`
    LastPrice   string  `json:"last_price"`
}

type Quote struct {
    Date    string  `json:"date"`
    Price   string     `json:"price"`
}

const (
    LayoutISO = "02/01/2006"
)
var assets []Asset
var DefaultDurations = [...]string {"1M","2M","3M","4M","5M","6M","7M","8M","9M","10M","11M","1Y","2Y","3Y"}
var DefaultPeriods = [...]string {"1","7","30","365"}
var quotes []Quote

func ScrapeSearchResult(query string) []Asset {
    doc := getHTMLDocument("https://www.boursorama.com/recherche/ajax?query=" + query)

    // Find the search results
	doc.Find(".search__list").First().Find(".search__list-link").Each(func(i int, s *goquery.Selection) {
        asset := Asset{}
		asset.Name = s.Find(".search__item-title").Text()
        link, ok := s.Attr("href")
        if !ok {
            log.Fatal("Unable to find link href")
        }
        splittedLink := strings.Split(link, "/")
        asset.Symbol = splittedLink[len(splittedLink)-2]
        asset.Category = strings.Trim(s.Find(".search__item-content").Text(), " \n")
        asset.LastPrice = s.Find(".search__item-instrument .last").Text()
        assets = append(assets, asset)
	})

    return assets 
}

func GetQuotes(symbol string, startDate time.Time, duration string, period string) []Quote {
    url := "https://www.boursorama.com/_formulaire-periode/?symbol=" + symbol + "&historic_search[startDate]=" + startDate.Format(LayoutISO) + "&historic_search[duration]=" + duration + "&historic_search[period]=" + period
    doc := getHTMLDocument(url)

    // Find the asset quotes
    doc.Find(".c-table tr").Each(func(i int, s *goquery.Selection) {
        // Escape first row (table header)
        if i == 0 {
            return
        }
        firstCell := s.Find(".c-table__cell").First()
        quote := Quote{}
        quote.Date = firstCell.Text()
        quote.Price = firstCell.Next().Text()
        quotes = append(quotes, quote)
    })

    return quotes
}

func getHTMLDocument(url string) *goquery.Document {
    // Request the HTML page
	res, err := http.Get(url)
	if err != nil {
        log.Fatal("Cannot fetch URL:", err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal("Cannot load HTML document", err)
	}
    return doc
}
