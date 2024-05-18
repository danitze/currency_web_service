package model

import (
	"genesis/currency-web-service/database"
	"gorm.io/gorm"
)

type Email struct {
	gorm.Model
	ID      uint
	Content string `gorm:"type:text" json:"content"`
}

func (email Email) Save() (*Email, error) {
	err := database.Database.Create(&email).Error
	if err != nil {
		return &Email{}, err
	}
	return &email, nil
}
