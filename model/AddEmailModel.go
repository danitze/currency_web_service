package model

type AddEmailModel struct {
	Email string `json:"email" validate:"required"`
}
