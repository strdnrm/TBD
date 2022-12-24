package model

import (
	"github.com/go-playground/validator"
	"github.com/google/uuid"
)

type Plane struct {
	Id            uuid.UUID `json:"-" db:"id"`
	NumberOfSeats int       `json:"number_of_seats" db:"number_of_seats" validate:"required,numeric"`
	Model         string    `json:"model" db:"model" validate:"required"`
}

func (p *Plane) Validate() error {
	validate := validator.New()
	err := validate.Struct(p)
	if err != nil {
		return err
	}
	return nil
}
