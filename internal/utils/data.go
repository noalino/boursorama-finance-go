package utils

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Asset struct {
	Symbol    string `json:"symbol"`
	LastPrice string `json:"last_price"`
	Market		string `json:"market"`
	Name      string `json:"name"`
}

type Quote struct {
	Date  string `json:"date"`
	Price string `json:"price"`
}

func ScrapeSearchResult(query string) ([]Asset, error) {
	url := getSearchUrl(query)
	doc, err := getHTMLDocument(url)
	if err != nil {
		return nil, err
	}

	// Find the searched results
	var assets []Asset
	view := doc.Find("[data-result-search]")

	view.Find("tbody.c-table__body").First().Find("tr.c-table__row").Each(func (i int, s *goquery.Selection) {
		asset := Asset{}
		cells := s.Find("td")
		link := cells.First().Find(".c-link")

		asset.Name = strings.TrimSpace(link.Text())

		url, ok := link.Attr("href")
		if !ok {
			log.Fatalf("Unable to find the quote symbol for %s\n", asset.Name)
			return
		}

		var symbolIndex int
		splitUrl := strings.Split(url, "/")
		runeUrl := []rune(url)
		if runeUrl[len(runeUrl)-1] == []rune("/")[0] {
			symbolIndex = len(splitUrl) - 2
		} else {
			symbolIndex = len(splitUrl) - 1
		}

		asset.Symbol = splitUrl[symbolIndex]
		asset.Market = strings.TrimSpace(cells.First().Next().Text())
		asset.LastPrice = strings.TrimSpace(cells.First().Next().Next().Text())
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

	// First page request to get the number of pages to scrape
	url := getQuotesUrl(symbol, startDate, duration, period, 1)
	doc, err := getHTMLDocument(url)
	if err != nil {
		return nil, err
	}

	nbOfPages := getMaxPages(doc)

	scrapeQuotes := func() []Quote {
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
		allQuotes = scrapeQuotes()
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
			return scrapeQuotes(), nil
		}
		// Init slice to return quotes from all pages
		quotesByPage := make([][]Quote, nbOfPages)
		// Use first page request to scrap quotes
		quotesByPage[0] = scrapeQuotes()
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
