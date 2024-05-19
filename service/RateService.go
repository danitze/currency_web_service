package service

import (
	"errors"
	"gorm.io/gorm"
	"time"
)

func GetRate(currency string, baseCurrency string) (float64, error) {
	rateModel, err := GetLocalRate(currency, baseCurrency)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			rate, err := FetchRate(currency, baseCurrency)
			if err != nil {
				return 0, err
			}

			rateModel, err := InsertRate(rate)
			if err != nil {
				return 0, err
			}
			return rateModel.Sale, nil
		} else {
			return 0, err
		}
	}
	if time.Now().Sub(rateModel.UpdateTime).Minutes() > 60 {
		rate, err := FetchRate(currency, baseCurrency)
		if err != nil {
			return 0, err
		}
		rateModel, err := UpdateRate(rate)
		if err != nil {
			return 0, err
		}
		return rateModel.Sale, nil
	} else {
		return rateModel.Sale, nil
	}
}
