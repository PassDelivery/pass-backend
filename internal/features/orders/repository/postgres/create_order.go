package orders_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	core_errors "github.com/PassDelivery/pass-backend/internal/core/errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *OrdersRepository) CreateOrder(
	ctx context.Context,
	userID uuid.UUID,
	order domain.CreateOrderRequest,
) (domain.OrderResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return domain.OrderResponse{}, fmt.Errorf("failed to begin: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil {
			return
		}
	}()

	queryAddress := `
	SELECT concat_ws(', ', city, street, house, 'кв. ' || apartment)
	FROM addresses
	WHERE id = $1 AND user_id = $2;
	`

	queryCart := `
	SELECT id, restaurant_id
	FROM carts
	WHERE user_id = $1
	FOR UPDATE;
	`

	queryCartItems := `
	SELECT ci.dish_id, d.name, d.price_kopecks, ci.quantity, d.is_available
	FROM cart_items ci
	JOIN dishes d ON d.id = ci.dish_id
	WHERE ci.cart_id = $1;
	`

	queryRestaurant := `
	SELECT name, is_open
	FROM restaurants
	WHERE id = $1;
	`

	queryOrder := `
	INSERT INTO orders (user_id, restaurant_id, status, total_kopecks, delivery_address, comment)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id, user_id, restaurant_id, status, total_kopecks, delivery_address, comment, created_at, updated_at;
	`

	queryItem := `
	INSERT INTO order_items (order_id, dish_id, dish_name, quantity, price_at_order_time)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id;
	`

	queryHistory := `
	INSERT INTO order_status_history (order_id, status, changed_by)
	VALUES ($1, $2, $3);
	`

	queryClearItems := `
	DELETE FROM cart_items
	WHERE cart_id = $1;
	`

	queryClearCart := `
	UPDATE carts
	SET restaurant_id = NULL, updated_at = NOW()
	WHERE id = $1;
	`

	var deliveryAddress string

	err = tx.QueryRow(ctx, queryAddress, order.AddressID, userID).Scan(&deliveryAddress)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.OrderResponse{}, fmt.Errorf("address not found: %w", core_errors.ErrNotFound)
	}
	if err != nil {
		return domain.OrderResponse{}, fmt.Errorf("scan address: %w", err)
	}

	var (
		cartID uuid.UUID
		restID *uuid.UUID
	)

	err = tx.QueryRow(ctx, queryCart, userID).Scan(&cartID, &restID)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.OrderResponse{}, fmt.Errorf("cart is empty: %w", core_errors.ErrConflict)
	}
	if err != nil {
		return domain.OrderResponse{}, fmt.Errorf("scan cart: %w", err)
	}

	rows, err := tx.Query(ctx, queryCartItems, cartID)
	if err != nil {
		return domain.OrderResponse{}, fmt.Errorf("query cart items: %w", err)
	}

	var (
		items []domain.OrderItem
		total int64
	)

	for rows.Next() {
		var (
			item        domain.OrderItem
			isAvailable bool
		)

		if err := rows.Scan(
			&item.DishID,
			&item.DishName,
			&item.PriceAtOrderTime,
			&item.Quantity,
			&isAvailable,
		); err != nil {
			rows.Close()
			return domain.OrderResponse{}, fmt.Errorf("scan cart item: %w", err)
		}
		if !isAvailable {
			rows.Close()
			return domain.OrderResponse{}, fmt.Errorf("dish %q is not available: %w", item.DishName, core_errors.ErrConflict)
		}

		total += item.PriceAtOrderTime * int64(item.Quantity)
		items = append(items, item)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return domain.OrderResponse{}, fmt.Errorf("cart items rows: %w", err)
	}

	if len(items) == 0 || restID == nil {
		return domain.OrderResponse{}, fmt.Errorf("cart is empty: %w", core_errors.ErrConflict)
	}

	var (
		restName string
		isOpen   bool
	)

	if err := tx.QueryRow(ctx, queryRestaurant, *restID).Scan(&restName, &isOpen); err != nil {
		return domain.OrderResponse{}, fmt.Errorf("scan restaurant: %w", err)
	}
	if !isOpen {
		return domain.OrderResponse{}, fmt.Errorf("restaurant is closed: %w", core_errors.ErrConflict)
	}

	var resp domain.OrderResponse

	if err := tx.QueryRow(
		ctx,
		queryOrder,
		userID,
		*restID,
		domain.OrderStatusNew,
		total,
		deliveryAddress,
		order.Comment,
	).Scan(
		&resp.ID,
		&resp.UserID,
		&resp.RestID,
		&resp.Status,
		&resp.TotalKopecks,
		&resp.DeliveryAddress,
		&resp.Comment,
		&resp.CreatedAt,
		&resp.UpdatedAt,
	); err != nil {
		return domain.OrderResponse{}, fmt.Errorf("insert order: %w", err)
	}

	for k := range items {
		if err := tx.QueryRow(
			ctx,
			queryItem,
			resp.ID,
			items[k].DishID,
			items[k].DishName,
			items[k].Quantity,
			items[k].PriceAtOrderTime,
		).Scan(&items[k].ID); err != nil {
			return domain.OrderResponse{}, fmt.Errorf("insert order item: %w", err)
		}
	}

	if _, err := tx.Exec(ctx, queryHistory, resp.ID, domain.OrderStatusNew, userID); err != nil {
		return domain.OrderResponse{}, fmt.Errorf("insert status history: %w", err)
	}

	if _, err := tx.Exec(ctx, queryClearItems, cartID); err != nil {
		return domain.OrderResponse{}, fmt.Errorf("clear cart items: %w", err)
	}

	if _, err := tx.Exec(ctx, queryClearCart, cartID); err != nil {
		return domain.OrderResponse{}, fmt.Errorf("reset cart: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return domain.OrderResponse{}, fmt.Errorf("faild commit to database: %w", err)
	}

	resp.RestName = restName
	resp.Items = items

	return resp, nil
}
