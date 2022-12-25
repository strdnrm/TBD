package sqlstore

import (
	"context"
	"final_project/internal/app/model"

	"github.com/google/uuid"
)

type RouteRepository struct {
	store *Store
}

func (r *RouteRepository) Create(ctx context.Context, rt *model.Route) error {
	if err := rt.Validate(); err != nil {
		return err
	}

	rt.Id = uuid.New()
	_, err := r.store.db.NamedQueryContext(ctx, `
	INSERT INTO route (id, source,  destination)
	VALUES (:id,:source, :destination);
	`, rt)
	if err != nil {
		return err
	}
	return nil
}
