package store

type Store interface {
	User() UserRepository
	Plane() PlaneRepository
}
