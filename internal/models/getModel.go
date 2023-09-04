package models

type Quote struct {
	Date  string  `json:"date"`
	Price float64 `json:"price"`
}

type Quotes []Quote
