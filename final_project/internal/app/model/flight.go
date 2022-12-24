package model

import (
	"github.com/go-playground/validator"
	"github.com/google/uuid"
)

type Flight struct {
	Id               uuid.UUID `json:"id" db:"id"`
	PlaneID          uuid.UUID `json:"plane_id" db:"plane_id" validate:"required"`
	RouteID          uuid.UUID `json:"route_id" db:"route_id" validate:"required"`
	DepartureTime    string    `json:"departure_time" db:"departure_time" validate:"required"`
	ArrivalTime      string    `json:"arrival_time" db:"arrival_time" validate:"required"`
	AvailableSeats   int       `json:"available_seats" db:"available_seats" validate:"required"`
	TransferFlightID uuid.UUID `json:"transfer_flight_id" db:"transfer_flight_id" validate:"omitempty"`
	Source           string    `json:"source" db:"source"`
	Destination      string    `json:"destintaion" db:"destintaion"`
}

func (p *Plane) Validate() error {
	validate := validator.New()
	err := validate.Struct(p)
	if err != nil {
		return err
	}
	return nil
}
