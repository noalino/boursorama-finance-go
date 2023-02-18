package utils

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Asset struct {
	Symbol    string `json:"symbol"`
	Name      string `json:"name"`
	LastPrice string `json:"last_price"`
}

type Quote struct {
	Date  string `json:"date"`
	Price string `json:"price"`
}

const (
	LayoutISO = "02/01/2006"
)

var DefaultDurations = []string{"1M", "2M", "3M", "4M", "5M", "6M", "7M", "8M", "9M", "10M", "11M", "1Y", "2Y", "3Y"}
var DefaultPeriods = []string{"1", "7", "30", "365"}

func ScrapeSearchResult(query string) ([]Asset, error) {
	sanitizedQuery := url.QueryEscape(strings.TrimSpace(query))
	doc, err := getHTMLDocument("https://www.boursorama.com/recherche/ajax?query=" + sanitizedQuery)
	if err != nil {
		return nil, err
	}

	// Find the search results
	var assets []Asset
	doc.Find(".search__list").First().Find(".search__list-link").Each(func(i int, s *goquery.Selection) {
		asset := Asset{}

		otherInfo := strings.Trim(s.Find(".search__item-content").Text(), " \n")
		name := s.Find(".search__item-title").Text()
		asset.Name = name + "\n" + otherInfo

		link, ok := s.Attr("href")
		if !ok {
			log.Fatalf("Unable to find the quote symbol for %s\n", asset.Name)
			return
		}
		var symbolIndex int
		splittedLink := strings.Split(link, "/")
		runeLink := []rune(link)
		if runeLink[len(runeLink)-1] == []rune("/")[0] {
			symbolIndex = len(splittedLink) - 2
		} else {
			symbolIndex = len(splittedLink) - 1
		}
		asset.Symbol = splittedLink[symbolIndex]

		asset.LastPrice = s.Find(".search__item-instrument .last").Text()

		assets = append(assets, asset)
	})

	return assets, nil
}

func GetQuotes(symbol string, startDate time.Time, duration string, period string) ([]Quote, error) {
	if ok := contains(DefaultDurations, duration); !ok {
		return nil, fmt.Errorf("Duration must be one of %v", DefaultDurations)
	}
	if ok := contains(DefaultPeriods, period); !ok {
		return nil, fmt.Errorf("Period must be one of %v", DefaultPeriods)
	}

	// First page request to get the number of pages to scrap
	url := getQuotesUrl(symbol, startDate, duration, period, 1)
	doc, err := getHTMLDocument(url)
	if err != nil {
		return nil, err
	}

	nbOfPages := doc.Find("span.c-pagination__content").Length()

	scrapQuotes := func() []Quote {
		quotes := []Quote{}
		doc.Find(".c-table tr").Each(func(i int, s *goquery.Selection) {
			// Escape first row (table header)
			if i == 0 {
				return
			}
			firstCell := s.Find(".c-table__cell").First()
			quote := Quote{}
			quote.Date = strings.TrimSpace(firstCell.Text())
			quote.Price = strings.TrimSpace(firstCell.Next().Text())
			quotes = append(quotes, quote)
		})
		return quotes
	}

	var allQuotes []Quote

	// Fetch quotes concurrently if there is more than one page
	if nbOfPages < 2 {
		allQuotes = scrapQuotes()
	} else {
		// Make channels to pass fatal errors in WaitGroup
		fatalErrors := make(chan error)
		wgDone := make(chan bool)

		var wg sync.WaitGroup
		// Scrap by page
		getPageQuotes := func(index int) ([]Quote, error) {
			url = getQuotesUrl(symbol, startDate, duration, period, index+1)
			doc, err = getHTMLDocument(url)
			if err != nil {
				return nil, err
			}
			return scrapQuotes(), nil
		}
		// Init slice to return quotes from all pages
		quotesByPage := make([][]Quote, nbOfPages)
		// Use first page request to scrap quotes
		quotesByPage[0] = scrapQuotes()
		// Fetch the remaining pages
		for i := 1; i < nbOfPages; i++ {
			wg.Add(1)

			go func(index int) {
				defer wg.Done()

				quotesByPage[index], err = getPageQuotes(index)
				if err != nil {
					fatalErrors <- err
				}
			}(i)
		}

		// Final goroutine to wait until WaitGroup is done
		go func() {
			wg.Wait()
			close(wgDone)
		}()

		// Wait until either WaitGroup is done or an error is received through the channel
		select {
		case <-wgDone:
			// Carry on
			break
		case err := <-fatalErrors:
			close(fatalErrors)
			return nil, err
		}

		for _, currentPageQuotes := range quotesByPage {
			allQuotes = append(allQuotes, currentPageQuotes...)
		}
	}

	return allQuotes, nil
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

func getQuotesUrl(symbol string, startDate time.Time, duration string, period string, page int) string {
	sanitizedSymbol := strings.ToUpper(strings.TrimSpace(symbol))
	if page == 1 {
		return "https://www.boursorama.com/_formulaire-periode/?symbol=" + sanitizedSymbol + "&historic_search[startDate]=" + startDate.Format(LayoutISO) + "&historic_search[duration]=" + duration + "&historic_search[period]=" + period
	} else {
		return "https://www.boursorama.com/_formulaire-periode/page-" + strconv.Itoa(page) + "?symbol=" + sanitizedSymbol + "&historic_search[startDate]=" + startDate.Format(LayoutISO) + "&historic_search[duration]=" + duration + "&historic_search[period]=" + period
	}
}

func contains(values []string, query string) bool {
	for _, value := range values {
		if value == query {
			return true
		}
	}
	return false
}
