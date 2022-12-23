package sqlstore

import (
	"context"
	"final_project/internal/app/model"
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

	_, err := r.store.db.NamedQueryContext(ctx, `
	INSERT INTO 
	users(login, email, password, name, surname, is_admin)
	VALUES (:login, :email, :password, :name, :surname, :is_admin)
	RETURNING ID; 
	`, u)
	//rw.Scan(&u.Id)
	if err != nil {
		return err
	}
	return nil
}
