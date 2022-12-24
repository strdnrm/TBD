package sqlstore

import (
	"final_project/internal/app/store"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Store struct {
	db              *sqlx.DB
	userRepository  *UserRepository
	planeRepository *PlaneRepository
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
