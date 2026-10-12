package cart_postgres_repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (r *CartRepository) ClearCart(
	ctx context.Context,
	userID uuid.UUID,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil {
			return
		}
	}()

	queryItems := `
	DELETE FROM cart_items ci
	USING carts c
	WHERE ci.cart_id = c.id AND c.user_id = $1;
	`

	queryCart := `
	UPDATE carts
	SET restaurant_id = NULL, updated_at = NOW()
	WHERE user_id = $1;
	`

	if _, err := tx.Exec(ctx, queryItems, userID); err != nil {
		return fmt.Errorf("delete cart items: %w", err)
	}

	if _, err := tx.Exec(ctx, queryCart, userID); err != nil {
		return fmt.Errorf("reset cart: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("faild commit to database: %w", err)
	}

	return nil
}
