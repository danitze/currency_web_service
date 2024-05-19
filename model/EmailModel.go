package model

import (
	"genesis/currency-web-service/apperror"
	"genesis/currency-web-service/database"
	"gorm.io/gorm"
)

type EmailModel struct {
	gorm.Model
	ID      uint
	Content string `json:"content"`
}

func (*EmailModel) TableName() string {
	return "emails"
}

func (email *EmailModel) Save() (*EmailModel, error) {
	var emails []EmailModel
	err := database.Database.Where("content = ?", email.Content).Find(&emails).Error
	if err != nil {
		return &EmailModel{}, err
	}
	if len(emails) > 0 {
		return &EmailModel{}, &apperror.DuplicateEmailError{}
	}
	err = database.Database.Create(email).Error
	if err != nil {
		return &EmailModel{}, err
	}
	return email, nil
}
