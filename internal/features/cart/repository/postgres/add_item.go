package cart_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	core_errors "github.com/PassDelivery/pass-backend/internal/core/errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *CartRepository) AddItem(
	ctx context.Context,
	userID uuid.UUID,
	item domain.CartItemRequest,
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

	queryDish := `
	SELECT restaurant_id, is_available
	FROM dishes
	WHERE id = $1;
	`

	queryCart := `
	INSERT INTO carts (user_id)
	VALUES ($1)
	ON CONFLICT (user_id) DO UPDATE SET updated_at = NOW()
	RETURNING id, restaurant_id;
	`

	queryHasItems := `
	SELECT EXISTS(SELECT 1 FROM cart_items WHERE cart_id = $1);
	`

	querySetRest := `
	UPDATE carts
	SET restaurant_id = $1
	WHERE id = $2;
	`

	queryItem := `
	INSERT INTO cart_items (cart_id, dish_id, quantity)
	VALUES ($1, $2, $3)
	ON CONFLICT (cart_id, dish_id) DO UPDATE SET quantity = cart_items.quantity + EXCLUDED.quantity;
	`

	var (
		dishRestID  uuid.UUID
		isAvailable bool
	)

	err = tx.QueryRow(ctx, queryDish, item.DishID).Scan(&dishRestID, &isAvailable)
	if errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf("dish not found: %w", core_errors.ErrNotFound)
	}
	if err != nil {
		return fmt.Errorf("scan dish: %w", err)
	}
	if !isAvailable {
		return fmt.Errorf("dish is not available: %w", core_errors.ErrConflict)
	}

	var (
		cartID     uuid.UUID
		cartRestID *uuid.UUID
	)

	if err := tx.QueryRow(ctx, queryCart, userID).Scan(&cartID, &cartRestID); err != nil {
		return fmt.Errorf("upsert cart: %w", err)
	}

	if cartRestID != nil && *cartRestID != dishRestID {
		var hasItems bool

		if err := tx.QueryRow(ctx, queryHasItems, cartID).Scan(&hasItems); err != nil {
			return fmt.Errorf("check cart items: %w", err)
		}
		if hasItems {
			return fmt.Errorf("cart contains dishes from another restaurant: %w", core_errors.ErrConflict)
		}
	}

	if _, err := tx.Exec(ctx, querySetRest, dishRestID, cartID); err != nil {
		return fmt.Errorf("set cart restaurant: %w", err)
	}

	if _, err := tx.Exec(ctx, queryItem, cartID, item.DishID, item.Quantity); err != nil {
		return fmt.Errorf("upsert cart item: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("faild commit to database: %w", err)
	}

	return nil
}
