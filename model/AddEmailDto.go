package model

type AddEmailDto struct {
	Email string `json:"email" validate:"required"`
}
