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

// Implement Sort interface
func (q Quotes) Len() int {
	return len(q)
}

func (q Quotes) Less(i, j int) bool {
	return q[i].Date.Before(q[j].Date.Time)
}

func (q Quotes) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
}
