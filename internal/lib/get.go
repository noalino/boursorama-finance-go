package lib

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"

	options "github.com/noalino/boursorama-finance-go/internal/lib/options/get"
	"github.com/noalino/boursorama-finance-go/internal/utils"
)

type GetQuery struct {
	Symbol   string
	From     string
	Duration string
	Period   string
}

type GetResultDate struct {
	time.Time
}

type GetResult struct {
	Date  GetResultDate `json:"date"`
	Price float64       `json:"price"`
}

type GetResults []GetResult

const GetResultDateFormat = "02/01/2006"

func (q GetQuery) Validate() (GetQuery, error) {
	q.Symbol = utils.SanitizeUrlInput(q.Symbol)
	if q.Symbol == "" {
		return q, errors.New("symbol value must be valid and not empty")
	}
	var err error
	q.From, err = options.From(q.From).ConvertToInternal()
	if err != nil {
		return q, err
	}
	q.Duration, err = options.Duration(q.Duration).ConvertToInternal()
	if err != nil {
		return q, err
	}
	q.Period, err = options.Period(q.Period).ConvertToInternal()
	if err != nil {
		return q, err
	}

	return q, nil
}

// Format date to JSON
func (d GetResultDate) MarshalJSON() ([]byte, error) {
	date := fmt.Sprintf("\"%s\"", d.Format(GetResultDateFormat))
	return []byte(date), nil
}

// Implement Sort interface
func (q GetResults) Len() int {
	return len(q)
}

func (q GetResults) Less(i, j int) bool {
	return q[i].Date.Before(q[j].Date.Time)
}

func (q GetResults) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
}

func getQuotesUrl(query GetQuery, page int) string {
	if page == 1 {
		return fmt.Sprintf(
			"%s/_formulaire-periode/?symbol=%s&historic_search[startDate]=%s&historic_search[duration]=%s&historic_search[period]=%s",
			utils.BASE_URL,
			strings.ToUpper(query.Symbol),
			query.From,
			query.Duration,
			query.Period,
		)
	} else {
		return fmt.Sprintf(
			"%s/_formulaire-periode/page-%s?symbol=%s&historic_search[startDate]=%s&historic_search[duration]=%s&historic_search[period]=%s",
			utils.BASE_URL,
			strconv.Itoa(page),
			strings.ToUpper(query.Symbol),
			query.From,
			query.Duration,
			query.Period,
		)
	}
}

func Get(unsafeQuery GetQuery) (GetResults, error) {
	query, err := unsafeQuery.Validate()
	if err != nil {
		return nil, err
	}

	// First page request to get the number of pages to scrape
	url := getQuotesUrl(query, 1)
	doc, err := utils.GetHTMLDocument(url)
	if err != nil {
		return nil, err
	}

	nbOfPages := utils.GetMaxPages(doc)

	scrapeQuotes := func() GetResults {
		quotes := GetResults{}
		doc.Find(".c-table tr").Each(func(i int, s *goquery.Selection) {
			// Escape first row (table header)
			if i == 0 {
				return
			}
			firstCell := s.Find(".c-table__cell").First()
			quote := GetResult{}
			date, err := time.Parse(GetResultDateFormat, strings.TrimSpace(firstCell.Text()))
			if err != nil {
				return
			}
			quote.Date = GetResultDate{Time: date}
			price, err := strconv.ParseFloat(strings.ReplaceAll(strings.TrimSpace(firstCell.Next().Text()), " ", ""), 64)
			if err != nil {
				quote.Price = 0.0
			}
			quote.Price = price
			quotes = append(quotes, quote)
		})
		return quotes
	}

	var allQuotes GetResults

	// Fetch quotes concurrently if there is more than one page
	if nbOfPages < 2 {
		allQuotes = scrapeQuotes()
	} else {
		// Make channels to pass fatal errors in WaitGroup
		fatalErrors := make(chan error)
		wgDone := make(chan bool)

		var wg sync.WaitGroup
		// Scrape by page
		getPageQuotes := func(index int) (GetResults, error) {
			url = getQuotesUrl(query, index+1)
			doc, err = utils.GetHTMLDocument(url)
			if err != nil {
				return nil, err
			}
			return scrapeQuotes(), nil
		}
		// Init slice to return quotes from all pages
		quotesByPage := make([]GetResults, nbOfPages)
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

	sort.Sort(allQuotes)
	return allQuotes, nil
}
