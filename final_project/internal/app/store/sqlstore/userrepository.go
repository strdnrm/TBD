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
	// _, err := r.store.db.QueryContext(ctx, `
	// INSERT INTO
	// users(id, login, email, password, name, surname, is_admin)
	// VALUES ($1, $2, $3, $4, $5, $6, $7)
	// RETURNING ID;
	// `, u.Id, u.Login, u.Email, u.EncryptedPassword, u.Name, u.Surname, u.Is_admin)
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

// func (r *UserRepository) FindByLogin(ctx context.Context, login string) (*model.User, error) {
// 	u := model.User{}
// 	err := r.store.db.GetContext(ctx, &u, `
// 	SELECT * FROM users WHERE login = $1;
// 	`, login)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &u, nil
// }
