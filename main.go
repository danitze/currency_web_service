package main

import (
	"encoding/json"
	"genesis/currency-web-service/database"
	"genesis/currency-web-service/model"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Rates []struct {
	Currency     string `json:"ccy"`
	BaseCurrency string `json:"base_ccy"`
	Buy          string `json:"buy"`
	Sale         string `json:"sale"`
}

type ResultCurrency struct {
	Buy  float64 `json:"buy"`
	Sale float64 `json:"sale"`
}

func main() {
	loadEnv()
	loadDatabase()
	router := gin.Default()
	router.GET("/currencies", getCurrencies)
	router.POST("/emails", addEmail)
	router.Run("localhost:8080")
}

func loadDatabase() {
	database.Connect()
	database.Database.AutoMigrate(&model.Email{})
}

func loadEnv() {
	err := godotenv.Load(".env.local")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func getCurrencies(c *gin.Context) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	req, err := http.NewRequest("GET", "https://api.privatbank.ua/p24api/pubinfo?json&exchange&coursid=5", nil)
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to make request: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	defer closeBody(resp.Body)
	if resp.StatusCode != http.StatusOK {
		log.Printf("Unexpected response status: %v", resp.StatusCode)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch currency data"})
		return
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	var rates Rates
	if err := json.Unmarshal(body, &rates); err != nil {
		log.Printf("Failed to unmarshal response body: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	var resultCurrency ResultCurrency
	for _, rate := range rates {
		if rate.Currency == "USD" {
			buy, err := strconv.ParseFloat(rate.Buy, 64)
			if err != nil {
				log.Printf("Cannot convert to float: %v", rate.Buy)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
				return
			}
			sale, err := strconv.ParseFloat(rate.Sale, 64)
			if err != nil {
				log.Printf("Cannot convert to float: %v", rate.Sale)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
				return
			}
			resultCurrency = ResultCurrency{
				Buy:  buy,
				Sale: sale,
			}
			c.IndentedJSON(http.StatusOK, resultCurrency)
			return
		}
	}
	c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
}

func addEmail(c *gin.Context) {
	validator := validator.New()
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot read body"})
		return
	}
	var addEmailModel model.AddEmailModel
	if err := json.Unmarshal(body, &addEmailModel); err != nil {
		log.Printf("Failed to unmarshal request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request body unmarshal error"})
		return
	}
	err = validator.Struct(&addEmailModel)
	if err != nil {
		log.Printf("Failed to validate request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request body validation error"})
		return
	}
	emailModel, err := model.Email.Save(model.Email{Content: addEmailModel.Email})
	if err != nil {
		log.Printf("Failed to insert email: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	c.IndentedJSON(http.StatusOK, emailModel)
}

func closeBody(body io.ReadCloser) {
	err := body.Close()
	if err != nil {
		log.Fatal(err)
	}
}