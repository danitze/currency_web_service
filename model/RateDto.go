package model

type RateDto struct {
	Currency     string  `json:"currency"`
	BaseCurrency string  `json:"base_currency"`
	Buy          float64 `json:"buy"`
	Sale         float64 `json:"sale"`
}
