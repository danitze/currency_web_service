package model

import (
	"genesis/currency-web-service/database"
	"gorm.io/gorm"
	"time"
)

type RateModel struct {
	gorm.Model
	ID           uint
	Currency     string    `json:"currency"`
	BaseCurrency string    `json:"base_currency"`
	Buy          float64   `json:"buy"`
	Sale         float64   `json:"sale"`
	UpdateTime   time.Time `json:"update_time"`
}

func (*RateModel) TableName() string {
	return "rates"
}

func (rate *RateModel) Save() (*RateModel, error) {
	err := database.Database.Create(rate).Error
	if err != nil {
		return nil, err
	}
	return rate, nil
}
