package main

import (
	"encoding/json"
	"errors"
	"genesis/currency-web-service/database"
	"genesis/currency-web-service/model"
	"genesis/currency-web-service/service"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
	"io"
	"log"
	"net/http"
	"time"
)

const (
	Currency     = "USD"
	BaseCurrency = "UAH"
)

func main() {
	loadEnv()
	loadDatabase()
	router := gin.Default()
	router.GET("/currencies", getCurrency)
	router.POST("/emails", addEmail)
	router.Run("localhost:8080")
}

func loadDatabase() {
	database.Connect()
	database.Database.AutoMigrate(&model.EmailModel{})
	database.Database.AutoMigrate(&model.RateModel{})
}

func loadEnv() {
	err := godotenv.Load(".env.local")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func getCurrency(c *gin.Context) {
	rateModel, err := service.GetRate(Currency, BaseCurrency)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			rate, err := service.FetchRate(Currency, BaseCurrency)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
				return
			}

			rateModel, err := service.InsertRate(rate)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
				return
			}
			c.IndentedJSON(http.StatusOK, rateModel.ToRateDto())
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		}
		return
	}
	if time.Now().Sub(rateModel.UpdateTime).Minutes() > 60 {
		rate, err := service.FetchRate(Currency, BaseCurrency)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}
		rateModel, err := service.UpdateRate(rate)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}
		c.IndentedJSON(http.StatusOK, rateModel.ToRateDto())
	} else {
		c.IndentedJSON(http.StatusOK, rateModel.ToRateDto())
	}
}

func addEmail(c *gin.Context) {
	structValidator := validator.New()
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot read body"})
		return
	}
	var addEmailModel model.AddEmailDto
	if err := json.Unmarshal(body, &addEmailModel); err != nil {
		log.Printf("Failed to unmarshal request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request body unmarshal error"})
		return
	}
	err = structValidator.Struct(&addEmailModel)
	if err != nil {
		log.Printf("Failed to validate request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request body validation error"})
		return
	}
	emailModel, err := (&model.EmailModel{Content: addEmailModel.Email}).Save()
	if err != nil {
		log.Printf("Failed to insert email: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	c.IndentedJSON(http.StatusOK, emailModel)
}
