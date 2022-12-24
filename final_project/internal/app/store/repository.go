package store

import (
	"context"
	"final_project/internal/app/model"

	"github.com/google/uuid"
)

type UserRepository interface {
	Create(context.Context, *model.User) error
	FindByEmail(context.Context, string) (*model.User, error)
	FindByID(context.Context, uuid.UUID) (*model.User, error)
}

type PlaneRepository interface {
	Create(context.Context, *model.Plane) error
}

type RouteRepository interface {
	Create(context.Context, *model.Route) error
}
