package model

import (
	"github.com/go-playground/validator"
	"github.com/google/uuid"
)

type Ticket struct {
	Id         uuid.UUID `json:"id" db:"id"`
	UserID     uuid.UUID `json:"user_id" db:"user_id" validate:"required"`
	FlightID   uuid.UUID `json:"flight_id" db:"flight_id" validate:"required"`
	Price      int       `json:"price" db:"price" validate:"required"`
	SeatID     uuid.UUID `json:"seat_id" db:"seat_id" validate:"required"`
	SeatNumber string    `json:"seat_number" db:"seat_number" validate:"required"`
}

func (t *Ticket) Validate() error {
	validate := validator.New()
	err := validate.Struct(t)
	if err != nil {
		return err
	}
	return nil
}
