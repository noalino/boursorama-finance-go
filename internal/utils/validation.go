package utils

import (
	"errors"
	"net/url"
	"strings"

	"github.com/noalino/boursorama-finance-go/internal/options"
)

type QuotesQuery struct {
	Symbol   string
	From     string
	Duration string
	Period   string
}

func ValidateInput(input string) string {
	return url.QueryEscape(strings.TrimSpace(input))
}

func ValidateQuotesQuery(query QuotesQuery) (QuotesQuery, error) {
	symbol := ValidateInput(query.Symbol)
	if symbol == "" {
		return QuotesQuery{}, errors.New("symbol value must be valid and not empty")
	}
	from, err := options.From(query.From).ConvertToInternal()
	if err != nil {
		return QuotesQuery{}, err
	}
	duration, err := options.Duration(query.Duration).ConvertToInternal()
	if err != nil {
		return QuotesQuery{}, err
	}
	period, err := options.Period(query.Period).ConvertToInternal()
	if err != nil {
		return QuotesQuery{}, err
	}

	return QuotesQuery{symbol, from, duration, period}, nil
}
