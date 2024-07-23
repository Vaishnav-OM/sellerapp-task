package models

type Product struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	Colour       string  `json:"colour"`
	Dimensions   string  `json:"dimensions"`
	Price        float64 `json:"price"`
	CurrencyUnit string  `json:"currencyUnit"`
}
