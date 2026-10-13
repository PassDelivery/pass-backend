package cart_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *CartRepository) GetCart(
	ctx context.Context,
	userID uuid.UUID,
) (domain.Cart, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	queryCart := `
	SELECT c.id, c.restaurant_id, r.name
	FROM carts c
	LEFT JOIN restaurants r ON r.id = c.restaurant_id
	WHERE c.user_id = $1;
	`

	queryItems := `
	SELECT ci.id, ci.dish_id, d.name, d.price_kopecks, ci.quantity, d.is_available
	FROM cart_items ci
	JOIN dishes d ON d.id = ci.dish_id
	WHERE ci.cart_id = $1
	ORDER BY d.name;
	`

	cart := domain.Cart{
		Items: []domain.CartItem{},
	}

	var cartID uuid.UUID

	err := r.pool.QueryRow(ctx, queryCart, userID).Scan(
		&cartID,
		&cart.RestID,
		&cart.RestName,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return cart, nil
	}
	if err != nil {
		return domain.Cart{}, fmt.Errorf("scan cart: %w", err)
	}

	rows, err := r.pool.Query(ctx, queryItems, cartID)
	if err != nil {
		return domain.Cart{}, fmt.Errorf("query cart items: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var item domain.CartItem

		if err := rows.Scan(
			&item.ID,
			&item.DishID,
			&item.Name,
			&item.PriceKopecks,
			&item.Quantity,
			&item.IsAvailable,
		); err != nil {
			return domain.Cart{}, fmt.Errorf("scan cart item: %w", err)
		}

		cart.Items = append(cart.Items, item)
	}
	if err := rows.Err(); err != nil {
		return domain.Cart{}, fmt.Errorf("cart items rows: %w", err)
	}

	return cart, nil
}
