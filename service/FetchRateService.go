package service

import (
	"encoding/json"
	"errors"
	"genesis/currency-web-service/model"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

func FetchRate(currencyName string, baseCurrencyName string) (*model.Rate, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	req, err := http.NewRequest("GET", "https://api.privatbank.ua/p24api/pubinfo?json&exchange&coursid=5", nil)
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to make request: %v", err)
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Printf("Unexpected response status: %v", resp.StatusCode)
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body: %v", err)
		return nil, err
	}
	var privatRates []model.PrivatRateDto
	if err := json.Unmarshal(body, &privatRates); err != nil {
		log.Printf("Failed to unmarshal response body: %v", err)
		return nil, err
	}
	var rate model.Rate
	for _, privatRate := range privatRates {
		if privatRate.Currency == currencyName && privatRate.BaseCurrency == baseCurrencyName {
			buy, err := strconv.ParseFloat(privatRate.Buy, 64)
			if err != nil {
				log.Printf("Cannot convert to float: %v", rate.Buy)
				return nil, err
			}
			sale, err := strconv.ParseFloat(privatRate.Sale, 64)
			if err != nil {
				log.Printf("Cannot convert to float: %v", rate.Sale)
				return nil, err
			}
			rate = model.Rate{
				Currency:     privatRate.Currency,
				BaseCurrency: privatRate.BaseCurrency,
				Buy:          buy,
				Sale:         sale,
			}
			return &rate, nil
		}
	}
	return nil, errors.New("cannot find rate")
}
