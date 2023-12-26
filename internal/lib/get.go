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
	Close float64       `json:"close"`
	Open  float64       `json:"open"`
	Perf  string        `json:"performance"`
	High  float64       `json:"high"`
	Low   float64       `json:"low"`
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

func getHistoricalUrl(query GetQuery, page uint16) string {
	return fmt.Sprintf(
		"%s/_formulaire-periode/page-%d?symbol=%s&historic_search[startDate]=%s&historic_search[duration]=%s&historic_search[period]=%s",
		utils.BASE_URL,
		page,
		strings.ToUpper(query.Symbol),
		query.From,
		query.Duration,
		query.Period,
	)
}

func Get(unsafeQuery GetQuery) (GetResults, error) {
	query, err := unsafeQuery.Validate()
	if err != nil {
		return nil, err
	}

	// First page request to get the number of pages to scrape
	url := getHistoricalUrl(query, 1)
	doc, err := utils.GetHTMLDocument(url)
	if err != nil {
		return nil, err
	}

	nbOfPages := utils.GetMaxPages(doc)

	scrapeHistorical := func() GetResults {
		results := GetResults{}
		doc.Find(".c-table tr").Each(func(i int, s *goquery.Selection) {
			// Skip first row (table header)
			if i == 0 {
				return
			}

			values := s.Find(".c-table__cell").Map(func(_ int, item *goquery.Selection) string {
				return strings.TrimSpace(item.Text())
			})

			result := GetResult{}
			date, err := time.Parse(GetResultDateFormat, values[0])
			if err != nil {
				return
			}
			result.Date = GetResultDate{Time: date}

			result.Perf = values[2]

			close, err := strconv.ParseFloat(values[1], 64)
			if err != nil {
				result.Close = 0.0
			}
			result.Close = close

			open, err := strconv.ParseFloat(values[5], 64)
			if err != nil {
				result.Open = 0.0
			}
			result.Open = open

			high, err := strconv.ParseFloat(values[3], 64)
			if err != nil {
				result.High = 0.0
			}
			result.High = high

			low, err := strconv.ParseFloat(values[4], 64)
			if err != nil {
				result.Low = 0.0
			}
			result.Low = low

			results = append(results, result)
		})
		return results
	}

	var allData GetResults

	// Fetch historical data concurrently if there is more than one page
	if nbOfPages < 2 {
		allData = scrapeHistorical()
	} else {
		// Make channels to pass fatal errors in WaitGroup
		fatalErrors := make(chan error)
		wgDone := make(chan bool)

		var wg sync.WaitGroup
		// Scrape by page
		getPageHistorical := func(index uint16) (GetResults, error) {
			url = getHistoricalUrl(query, index+1)
			doc, err = utils.GetHTMLDocument(url)
			if err != nil {
				return nil, err
			}
			return scrapeHistorical(), nil
		}
		// Init slice to return data from all pages
		resultsByPage := make([]GetResults, nbOfPages)
		// Use first page request to scrape data
		resultsByPage[0] = scrapeHistorical()
		// Fetch the remaining pages
		for i := uint16(1); i < nbOfPages; i++ {
			wg.Add(1)

			go func(index uint16) {
				defer wg.Done()

				resultsByPage[index], err = getPageHistorical(index)
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

		for _, currentPageResults := range resultsByPage {
			allData = append(allData, currentPageResults...)
		}
	}

	sort.Sort(allData)
	return allData, nil
}
