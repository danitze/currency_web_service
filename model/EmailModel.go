package model

import (
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
	err := database.Database.Create(email).Error
	if err != nil {
		return &EmailModel{}, err
	}
	return email, nil
}
