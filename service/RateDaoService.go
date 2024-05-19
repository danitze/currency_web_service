package service

import (
	"genesis/currency-web-service/database"
	"genesis/currency-web-service/model"
	"log"
	"time"
)

func InsertRate(rate *model.Rate) (*model.RateModel, error) {
	rateModel, err := (&model.RateModel{
		Currency:     rate.Currency,
		BaseCurrency: rate.BaseCurrency,
		Buy:          rate.Buy,
		Sale:         rate.Sale,
		UpdateTime:   time.Now(),
	}).Save()
	if err != nil {
		log.Printf("Failed to insert rate: %v", err)
		return nil, err
	}
	return rateModel, err
}

func UpdateRate(rate *model.Rate) (*model.RateModel, error) {
	updateStatement := "UPDATE rates SET buy=$1, sale=$2, update_time=$3 WHERE currency=$4 AND base_currency=$5"
	err := database.Database.Exec(
		updateStatement, rate.Buy, rate.Sale, time.Now(), rate.Currency, rate.BaseCurrency,
	).Error
	if err != nil {
		log.Printf("Failed to update rate: %v", err)
		return nil, err
	}
	return GetRate(rate.Currency, rate.BaseCurrency)
}

func GetRate(currency string, baseCurrency string) (*model.RateModel, error) {
	var rateModel model.RateModel
	err := database.Database.Where("currency = ? AND base_currency = ?",
		currency, baseCurrency).First(&rateModel).Error
	return &rateModel, err
}
