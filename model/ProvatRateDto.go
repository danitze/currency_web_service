package model

type PrivatRateDto struct {
	Currency     string `json:"ccy"`
	BaseCurrency string `json:"base_ccy"`
	Buy          string `json:"buy"`
	Sale         string `json:"sale"`
}
