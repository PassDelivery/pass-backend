package menegment_postgres_repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (r *MenegmentRepository) IsDishInActiveOrders(
	ctx context.Context,
	dishID uuid.UUID,
) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT EXISTS(
		SELECT 1 FROM order_items oi
		JOIN orders o ON o.id = oi.order_id
		WHERE oi.dish_id = $1 AND o.status IN ('new', 'accepted', 'cooking', 'ready', 'delivering')
	);
	`

	var ok bool

	if err := r.pool.QueryRow(ctx, query, dishID).Scan(&ok); err != nil {
		return false, fmt.Errorf("failed scan: %w", err)
	}

	return ok, nil
}

func (r *MenegmentRepository) HideDish(
	ctx context.Context,
	dishID uuid.UUID,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	UPDATE dishes
	SET is_available = false
	WHERE id = $1;
	`

	if _, err := r.pool.Exec(ctx, query, dishID); err != nil {
		return fmt.Errorf("failed exec: %w", err)
	}

	return nil
}
