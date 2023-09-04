package utils

import (
	"log"
	"strconv"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"

	"github.com/noalino/boursorama-finance-go/internal/models"
)

func ScrapeSearchResult(unsafeQuery SearchQuery) (models.Assets, error) {
	query, err := unsafeQuery.Validate()
	if err != nil {
		return nil, err
	}

	url := getSearchUrl(query.Value)
	doc, err := getHTMLDocument(url)
	if err != nil {
		return nil, err
	}

	// Find the searched results
	var assets models.Assets
	view := doc.Find("[data-result-search]")

	view.Find("tbody.c-table__body").First().Find("tr.c-table__row").Each(func(i int, s *goquery.Selection) {
		asset := models.Asset{}
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

func GetQuotes(unsafeQuery QuotesQuery) (models.Quotes, error) {
	query, err := unsafeQuery.Validate()
	if err != nil {
		return nil, err
	}

	// First page request to get the number of pages to scrape
	url := getQuotesUrl(query, 1)
	doc, err := getHTMLDocument(url)
	if err != nil {
		return nil, err
	}

	nbOfPages := getMaxPages(doc)

	scrapeQuotes := func() models.Quotes {
		quotes := models.Quotes{}
		doc.Find(".c-table tr").Each(func(i int, s *goquery.Selection) {
			// Escape first row (table header)
			if i == 0 {
				return
			}
			firstCell := s.Find(".c-table__cell").First()
			quote := models.Quote{}
			quote.Date = strings.TrimSpace(firstCell.Text())
			price, err := strconv.ParseFloat(strings.ReplaceAll(strings.TrimSpace(firstCell.Next().Text()), " ", ""), 64)
			if err != nil {
				quote.Price = 0.0
			}
			quote.Price = price
			quotes = append(quotes, quote)
		})
		return quotes
	}

	var allQuotes models.Quotes

	// Fetch quotes concurrently if there is more than one page
	if nbOfPages < 2 {
		allQuotes = scrapeQuotes()
	} else {
		// Make channels to pass fatal errors in WaitGroup
		fatalErrors := make(chan error)
		wgDone := make(chan bool)

		var wg sync.WaitGroup
		// Scrape by page
		getPageQuotes := func(index int) (models.Quotes, error) {
			url = getQuotesUrl(query, index+1)
			doc, err = getHTMLDocument(url)
			if err != nil {
				return nil, err
			}
			return scrapeQuotes(), nil
		}
		// Init slice to return quotes from all pages
		quotesByPage := make([]models.Quotes, nbOfPages)
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
