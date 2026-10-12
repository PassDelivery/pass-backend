package cart_postgres_repository

import (
	"context"
	"fmt"

	core_errors "github.com/PassDelivery/pass-backend/internal/core/errors"
	"github.com/google/uuid"
)

func (r *CartRepository) PatchItem(
	ctx context.Context,
	userID uuid.UUID,
	itemID uuid.UUID,
	quantity int,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	UPDATE cart_items ci
	SET quantity = $1
	FROM carts c
	WHERE ci.id = $2 AND ci.cart_id = c.id AND c.user_id = $3;
	`

	tag, err := r.pool.Exec(ctx, query, quantity, itemID, userID)
	if err != nil {
		return fmt.Errorf("failed exec: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("cart item not found: %w", core_errors.ErrNotFound)
	}

	return nil
}
