package service

import (
	"genesis/currency-web-service/database"
	"genesis/currency-web-service/model"
)

func GetEmails() (*[]model.EmailModel, error) {
	var emails []model.EmailModel
	err := database.Database.Find(&emails).Error
	return &emails, err
}
