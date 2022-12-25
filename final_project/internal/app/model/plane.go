package model

import (
	"github.com/go-playground/validator"
	"github.com/google/uuid"
)

type Plane struct {
	Id            uuid.UUID `json:"id" db:"id"`
	NumberOfSeats int       `json:"number_of_seats" db:"number_of_seats" validate:"required,numeric"`
	Model         string    `json:"model" db:"model" validate:"required"`
}

func (f *Flight) Validate() error {
	validate := validator.New()
	err := validate.Struct(f)
	if err != nil {
		return err
	}
	return nil
}
