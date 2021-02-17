package utils 

import (
    "fmt"
    "log"
    "net/url"
    "net/http"
    "time"
    "strconv"
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
    Price   string  `json:"price"`
}

const (
    LayoutISO = "02/01/2006"
)
var DefaultDurations = []string{"1M","2M","3M","4M","5M","6M","7M","8M","9M","10M","11M","1Y","2Y","3Y"}
var DefaultPeriods = []string{"1","7","30","365"}

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

    return assets, nil
}

func GetQuotes(symbol string, startDate time.Time, duration string, period string) ([]Quote, error) {
    if ok := contains(DefaultDurations, duration); !ok {
        return nil, fmt.Errorf("Duration must be one of %v", DefaultDurations)
    }
    if ok := contains(DefaultPeriods, period); !ok {
        return nil, fmt.Errorf("Period must be one of %v", DefaultPeriods)
    }

    var quotes []Quote
    page := 1
    url := getQuotesUrl(symbol, startDate, duration, period, page)
    doc, err := getHTMLDocument(url)
    if err != nil {
        return nil, err
    }

    nbOfPages := doc.Find("span.c-pagination__content").Length()
    log.Printf("Number of pages: %d", nbOfPages)

    // Find the asset quotes
    appendQuotes := func() {
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
    }

    // Fetch all pages if any
    var getQuotes func() (bool, error)
    getQuotes = func() (bool, error) {
        log.Printf("URL: %s", url)
        log.Printf("page: %d", page)
        doc, err = getHTMLDocument(url)
        if err != nil {
            return false, err
        }
        appendQuotes()
        page = page + 1
        if nbOfPages != 0 && page <= nbOfPages {
            url = getQuotesUrl(symbol, startDate, duration, period, page)
            ok, err := getQuotes()
            if !ok {
                return false, err
            }
        }
        return true, nil
    } 

    ok, err := getQuotes()
    if !ok {
        return nil, err
    }

    return quotes, nil
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
