package sqlstore

import (
	"context"
	"final_project/internal/app/model"

	"github.com/google/uuid"
)

type FlightRepository struct {
	store *Store
}

func (r *FlightRepository) Create(ctx context.Context, f *model.Flight) error {
	if err := f.Validate(); err != nil {
		return err
	}

	f.Id = uuid.New()
	_, err := r.store.db.NamedQueryContext(ctx, `
	INSERT INTO
	flight(id, plane_id, route_id, departure_time, arrival_time, available_seats)
	VALUES (:id, :plane_id, :route_id, :departure_time, :arrival_time, :available_seats);
	`, f)
	if err != nil {
		return err
	}
	return nil
}

func (r *FlightRepository) GetByArrivalTime(ctx context.Context, arrivalTime string, rt *model.Route) (*model.Flight, error) {
	f := model.Flight{}
	err := r.store.db.GetContext(ctx, &f, `
	SELECT * FROM flight
	INNER JOIN route ON route.id = flight.id
	WHERE flight.arrival_time::date = $1::date
	AND route.source = $2
	AND route.destination = $3;
	`, arrivalTime, rt.Source, rt.Destination)
	if err != nil {
		return nil, err
	}
	return &f, nil
}

func (r *FlightRepository) GetByDepartureTime(ctx context.Context, arrivalTime string, rt *model.Route) (*model.Flight, error) {
	f := model.Flight{}
	err := r.store.db.GetContext(ctx, &f, `
	SELECT * FROM flight
	INNER JOIN route ON route.id = flight.id
	WHERE flight.departure_time::date = $1::date
	AND route.source = $2
	AND route.destination = $3;
	`, arrivalTime, rt.Source, rt.Destination)
	if err != nil {
		return nil, err
	}
	return &f, nil
}
