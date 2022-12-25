package sqlstore

import (
	"context"
	"final_project/internal/app/model"

	"github.com/google/uuid"
)

type UserRepository struct {
	store *Store
}

func (r *UserRepository) Create(ctx context.Context, u *model.User) error {
	if err := u.Validate(); err != nil {
		return err
	}

	if err := u.BeforeCreate(); err != nil {
		return err
	}
	u.Id = uuid.New()
	_, err := r.store.db.NamedQueryContext(ctx, `
	INSERT INTO
	users(id, login, email, password, name, surname, is_admin)
	VALUES (:id, :login, :email, :password, :name, :surname, :is_admin)
	RETURNING ID;
	`, u)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	u := model.User{}
	err := r.store.db.GetContext(ctx, &u, `
	SELECT * FROM users WHERE email = $1;
	`, email)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	u := model.User{}
	err := r.store.db.GetContext(ctx, &u, `
	SELECT * FROM users WHERE id = $1;
	`, id)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) FindByLogin(ctx context.Context, login string) (*model.User, error) {
	u := model.User{}
	err := r.store.db.GetContext(ctx, &u, `
	SELECT * FROM users WHERE login = $1;
	`, login)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) GetFlightsByDeparturePoint(ctx context.Context, departutePoint string, u *model.User) (*[]model.Flight, error) {
	f := []model.Flight{}
	err := r.store.db.SelectContext(ctx, &f, `
	SELECT * 
	FROM user_tickets
	LEFT JOIN ticket ON ticket.id = user_tickets.ticket_id
	LEFT JOIN flight ON flight.id = ticket.flight_id
	LEFT JOIN route ON route.id = flight.route_id
	WHERE user_tickets.user_id = $1
	AND route.destination = $2
	ORDER BY flight.departure_time DESC;
	`, u.Id, departutePoint)
	if err != nil {
		return nil, err
	}
	return &f, nil
}

func (r *UserRepository) GetFlightsByArrivalPoint(ctx context.Context, arrivalPoint string, u *model.User) (*[]model.Flight, error) {
	f := []model.Flight{}
	err := r.store.db.SelectContext(ctx, &f, `
	SELECT * 
	FROM user_tickets
	LEFT JOIN ticket ON ticket.id = user_tickets.ticket_id
	LEFT JOIN flight ON flight.id = ticket.flight_id
	LEFT JOIN route ON route.id = flight.route_id
	WHERE user_tickets.user_id = $1
	AND route.destination = $2
	ORDER BY flight.departure_time DESC;
	`, u.Id, arrivalPoint)
	if err != nil {
		return nil, err
	}
	return &f, nil
}

func (r *UserRepository) GetFlightsByDepartureDate(ctx context.Context, departuteDate string, u *model.User) (*[]model.Flight, error) {
	f := []model.Flight{}
	err := r.store.db.SelectContext(ctx, &f, `
	SELECT * 
	FROM user_tickets
	LEFT JOIN ticket ON ticket.id = user_tickets.ticket_id
	LEFT JOIN flight ON flight.id = ticket.flight_id
	LEFT JOIN route ON route.id = flight.route_id
	WHERE user_tickets.user_id = $1
	AND flight.departure_time::date = $2::date
	ORDER BY flight.departure_time DESC;
	`, u.Id, departuteDate)
	if err != nil {
		return nil, err
	}
	return &f, nil
}

func (r *UserRepository) GetFlightsByArrivalDate(ctx context.Context, arrivalDate string, u *model.User) (*[]model.Flight, error) {
	f := []model.Flight{}
	err := r.store.db.SelectContext(ctx, &f, `
	SELECT * 
	FROM user_tickets
	LEFT JOIN ticket ON ticket.id = user_tickets.ticket_id
	LEFT JOIN flight ON flight.id = ticket.flight_id
	LEFT JOIN route ON route.id = flight.route_id
	WHERE user_tickets.user_id = 1
	AND flight.arrival_time::date = '2020.04.04'::date
	ORDER BY flight.departure_time DESC;
	`, u.Id, arrivalDate)
	if err != nil {
		return nil, err
	}
	return &f, nil
}
