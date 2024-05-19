package main

import (
	"encoding/json"
	"fmt"
	"genesis/currency-web-service/apperror"
	"genesis/currency-web-service/database"
	"genesis/currency-web-service/model"
	"genesis/currency-web-service/service"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/robfig/cron"
	"io"
	"log"
	"net/http"
)

const (
	Currency     = "USD"
	BaseCurrency = "UAH"
)

func main() {
	loadEnv()
	loadDatabase()
	initEmailSending()
	router := gin.Default()
	router.GET("/rate", getRate)
	router.POST("/subscribe", subscribe)
	err := router.Run("localhost:8080")
	if err != nil {
		log.Fatal("Error running router")
	}
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

func initEmailSending() {
	appCron := cron.New()
	err := appCron.AddFunc("@midnight", func() {
		err := sendCurrencyToEmails()
		if err != nil {
			log.Printf("Failed to send current rate to emails: %v", err)
		}
	})
	if err != nil {
		log.Printf("Failed to init email sending: %v", err)
		return
	}
	appCron.Start()
}

func getRate(c *gin.Context) {
	rate, err := service.GetRate(Currency, BaseCurrency)
	if err != nil {
		log.Printf("Failed to get current rate: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	c.JSON(http.StatusOK, rate)
}

func subscribe(c *gin.Context) {
	structValidator := validator.New()
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"apperror": "Cannot read body"})
		return
	}
	var addEmailModel model.AddEmailDto
	if err := json.Unmarshal(body, &addEmailModel); err != nil {
		log.Printf("Failed to unmarshal request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"apperror": "Request body unmarshal apperror"})
		return
	}
	err = structValidator.Struct(&addEmailModel)
	if err != nil {
		log.Printf("Failed to validate request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"apperror": "Request body validation apperror"})
		return
	}
	_, err = (&model.EmailModel{Content: addEmailModel.Email}).Save()
	if err != nil {
		if err.Error() == apperror.DuplicateEmailErrorText {
			log.Printf("Duplicate email")
			c.JSON(http.StatusConflict, "Даний E-mail вже доданий до розсилки")
		} else {
			log.Printf("Failed to insert email: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"apperror": "Internal Server Error"})
		}
		return
	}
	c.IndentedJSON(http.StatusOK, "E-mail додано")
}

func sendCurrencyToEmails() error {
	emailModels, err := service.GetEmails()
	if err != nil {
		return err
	}
	emails := make([]string, len(*emailModels))
	for index, emailModel := range *emailModels {
		emails[index] = emailModel.Content
	}
	rate, err := service.GetRate(Currency, BaseCurrency)
	if err != nil {
		return err
	}
	err = service.SendToEmails(emails, "Поточний курс долара", fmt.Sprintf("Поточний курс долара до гривні: %v", rate))
	return err
}
