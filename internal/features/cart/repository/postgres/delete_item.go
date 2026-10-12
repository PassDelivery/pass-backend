package cart_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	core_errors "github.com/PassDelivery/pass-backend/internal/core/errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *CartRepository) DeleteItem(
	ctx context.Context,
	userID uuid.UUID,
	itemID uuid.UUID,
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

	queryDelete := `
	DELETE FROM cart_items ci
	USING carts c
	WHERE ci.id = $1 AND ci.cart_id = c.id AND c.user_id = $2
	RETURNING ci.cart_id;
	`

	queryResetRest := `
	UPDATE carts
	SET restaurant_id = NULL, updated_at = NOW()
	WHERE id = $1 AND NOT EXISTS(SELECT 1 FROM cart_items WHERE cart_id = $1);
	`

	var cartID uuid.UUID

	err = tx.QueryRow(ctx, queryDelete, itemID, userID).Scan(&cartID)
	if errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf("cart item not found: %w", core_errors.ErrNotFound)
	}
	if err != nil {
		return fmt.Errorf("delete cart item: %w", err)
	}

	if _, err := tx.Exec(ctx, queryResetRest, cartID); err != nil {
		return fmt.Errorf("reset cart restaurant: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("faild commit to database: %w", err)
	}

	return nil
}
