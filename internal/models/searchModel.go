package models

type Asset struct {
	Symbol    string `json:"symbol"`
	LastPrice string `json:"last_price"`
	Market    string `json:"market"`
	Name      string `json:"name"`
}

type Assets []Asset
