package sqlstore

import (
	"context"
	"final_project/internal/app/model"

	"github.com/google/uuid"
)

type TicketRepository struct {
	store *Store
}

func (r *TicketRepository) Purchase(ctx context.Context, t *model.Ticket) error {
	if err := t.Validate(); err != nil {
		return err
	}

	t.Id = uuid.New()
	_, err := r.store.db.NamedQueryContext(ctx, `
	INSERT INTO ticket (id, user_id, flight_id, price, seat_id)
	VALUES (:id, :user_id, :flight_id, :price, :seat_id);
	`, t)
	if err != nil {
		return err
	}

	_, err = r.store.db.QueryContext(ctx, `
	INSERT INTO user_tickets(user_id, ticket_id)
	VALUES($1, $2);
	`, t.UserID, t.Id)
	if err != nil {
		return err
	}

	_, err = r.store.db.ExecContext(ctx, `
	UPDATE seat
	SET user_id = $1
	WHERE seat_number = $2;
	`, t.UserID, t.SeatNumber)
	if err != nil {
		return err
	}

	return nil
}
