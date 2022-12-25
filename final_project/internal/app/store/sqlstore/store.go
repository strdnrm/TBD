package sqlstore

import (
	"final_project/internal/app/store"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Store struct {
	db               *sqlx.DB
	userRepository   *UserRepository
	planeRepository  *PlaneRepository
	routeRepository  *RouteRepository
	flightRepository *FlightRepository
	ticketRepository *TicketRepository
}

func New(db *sqlx.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) User() store.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		store: s,
	}

	return s.userRepository
}

func (s *Store) Plane() store.PlaneRepository {
	if s.planeRepository != nil {
		return s.planeRepository
	}

	s.planeRepository = &PlaneRepository{
		store: s,
	}

	return s.planeRepository
}

func (s *Store) Route() store.RouteRepository {
	if s.routeRepository != nil {
		return s.routeRepository
	}

	s.routeRepository = &RouteRepository{
		store: s,
	}

	return s.routeRepository
}

func (s *Store) Flight() store.FlightRepository {
	if s.flightRepository != nil {
		return s.flightRepository
	}

	s.flightRepository = &FlightRepository{
		store: s,
	}

	return s.flightRepository
}

func (s *Store) Ticket() store.TicketRepository {
	if s.ticketRepository != nil {
		return s.ticketRepository
	}

	s.ticketRepository = &TicketRepository{
		store: s,
	}

	return s.ticketRepository
}
