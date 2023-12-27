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
	Page  uint16
}

type SearchResultAsset struct {
	Symbol    string `json:"symbol"`
	LastPrice string `json:"last_price"`
	Market    string `json:"market"`
	Name      string `json:"name"`
}

type SearchResult struct {
	Page       uint16              `json:"current_page"`
	TotalPages uint16              `json:"total_pages_count"`
	Assets     []SearchResultAsset `json:"values"`
}

func (q SearchQuery) Validate() (SearchQuery, error) {
	q.Value = utils.SanitizeUrlInput(q.Value)
	if q.Value == "" {
		return q, errors.New("search value must be valid and not empty")
	}

	return q, nil
}

func getSearchUrl(query SearchQuery) string {
	return fmt.Sprintf("%s/recherche/_instruments/%s?page=%d", utils.BASE_URL, query.Value, query.Page)
}

func Search(unsafeQuery SearchQuery) (SearchResult, error) {
	query, err := unsafeQuery.Validate()
	if err != nil {
		return SearchResult{}, err
	}

	url := getSearchUrl(query)
	doc, err := utils.GetHTMLDocument(url)
	if err != nil {
		return SearchResult{}, err
	}

	nbOfPages := utils.GetMaxPages(doc)

	result := SearchResult{
		Page:       query.Page,
		TotalPages: nbOfPages,
		Assets:     []SearchResultAsset{},
	}

	// Find the searched results
	view := doc.Find("[data-result-search]")

	view.Find("tbody.c-table__body").First().Find("tr.c-table__row").Each(func(i int, s *goquery.Selection) {
		asset := SearchResultAsset{}
		cells := s.Find("td")
		link := cells.First().Find(".c-link")

		asset.Name = strings.TrimSpace(link.Text())

		url, ok := link.Attr("href")
		if !ok {
			log.Fatalf("Unable to find the symbol for %s\n", asset.Name)
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
		result.Assets = append(result.Assets, asset)
	})

	return result, nil
}
