package store

import (
	"context"
	"final_project/internal/app/model"
)

type UserRepository interface {
	Create(context.Context, *model.User) error
	FindByEmail(context.Context, string) (*model.User, error)
}
