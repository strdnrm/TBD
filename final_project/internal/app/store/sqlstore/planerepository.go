package sqlstore

import (
	"context"
	"final_project/internal/app/model"

	"github.com/google/uuid"
)

type PlaneRepository struct {
	store *Store
}

func (r *PlaneRepository) Create(ctx context.Context, p *model.Plane) error {
	if err := p.Validate(); err != nil {
		return err
	}

	p.Id = uuid.New()
	_, err := r.store.db.NamedQueryContext(ctx, `
	INSERT INTO
	plane(id, number_of_seats, model)
	VALUES (:id, :number_of_seats, :model)
	RETURNING ID;
	`, p)
	if err != nil {
		return err
	}
	return nil
}
