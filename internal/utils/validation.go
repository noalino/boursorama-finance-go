package utils

import (
	"errors"
	"net/url"
	"strings"

	"github.com/noalino/boursorama-finance-go/internal/options"
)

type SearchQuery struct {
	Value string
}

type QuotesQuery struct {
	Symbol   string
	From     string
	Duration string
	Period   string
}

func SanitizeUrlInput(input string) string {
	return url.QueryEscape(strings.TrimSpace(input))
}

func (q SearchQuery) Validate() (SearchQuery, error) {
	value := SanitizeUrlInput(q.Value)
	if value == "" {
		return SearchQuery{}, errors.New("search value must be valid and not empty")
	}

	return SearchQuery{value}, nil
}

func (q QuotesQuery) Validate() (QuotesQuery, error) {
	symbol := SanitizeUrlInput(q.Symbol)
	if symbol == "" {
		return QuotesQuery{}, errors.New("symbol value must be valid and not empty")
	}
	from, err := options.From(q.From).ConvertToInternal()
	if err != nil {
		return QuotesQuery{}, err
	}
	duration, err := options.Duration(q.Duration).ConvertToInternal()
	if err != nil {
		return QuotesQuery{}, err
	}
	period, err := options.Period(q.Period).ConvertToInternal()
	if err != nil {
		return QuotesQuery{}, err
	}

	return QuotesQuery{symbol, from, duration, period}, nil
}
