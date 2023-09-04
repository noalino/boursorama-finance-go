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
	q.Value = SanitizeUrlInput(q.Value)
	if q.Value == "" {
		return q, errors.New("search value must be valid and not empty")
	}

	return q, nil
}

func (q QuotesQuery) Validate() (QuotesQuery, error) {
	q.Symbol = SanitizeUrlInput(q.Symbol)
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
