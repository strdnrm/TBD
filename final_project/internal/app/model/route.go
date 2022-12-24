package model

import (
	"github.com/go-playground/validator"
	"github.com/google/uuid"
)

type Route struct {
	Id          uuid.UUID `json:"id" db:"id"`
	Source      string    `json:"source" db:"source" validate:"required"`
	Destination string    `json:"destintaion" db:"destintaion" validate:"required"`
}

func (r *Route) Validate() error {
	validate := validator.New()
	err := validate.Struct(r)
	if err != nil {
		return err
	}
	return nil
}
