package lib

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"github.com/noalino/boursorama-finance-go/internal/utils"
)

type SearchQuery struct {
	Value string
}

type SearchResult struct {
	Symbol    string `json:"symbol"`
	LastPrice string `json:"last_price"`
	Market    string `json:"market"`
	Name      string `json:"name"`
}

type SearchResults []SearchResult

func (q SearchQuery) Validate() (SearchQuery, error) {
	q.Value = utils.SanitizeUrlInput(q.Value)
	if q.Value == "" {
		return q, errors.New("search value must be valid and not empty")
	}

	return q, nil
}

func getSearchUrl(searchValue string) string {
	return fmt.Sprintf("%s/recherche/_instruments/%s", utils.BASE_URL, searchValue)
}

func Search(unsafeQuery SearchQuery) (SearchResults, error) {
	query, err := unsafeQuery.Validate()
	if err != nil {
		return nil, err
	}

	url := getSearchUrl(query.Value)
	doc, err := utils.GetHTMLDocument(url)
	if err != nil {
		return nil, err
	}

	// Find the searched results
	var assets SearchResults
	view := doc.Find("[data-result-search]")

	view.Find("tbody.c-table__body").First().Find("tr.c-table__row").Each(func(i int, s *goquery.Selection) {
		asset := SearchResult{}
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
