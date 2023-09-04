package models

import (
	"fmt"
	"time"
)

type QuoteDate struct {
	time.Time
}

type Quote struct {
	Date  QuoteDate `json:"date"`
	Price float64   `json:"price"`
}

type Quotes []Quote

const DateFormat = "02/01/2006"

// Format date to JSON
func (d QuoteDate) MarshalJSON() ([]byte, error) {
	date := fmt.Sprintf("\"%s\"", d.Format(DateFormat))
	return []byte(date), nil
}
