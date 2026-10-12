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

func (r *OrdersRepository) GetOrder(
	ctx context.Context,
	orderID uuid.UUID,
) (domain.OrderResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT o.id, o.user_id, o.restaurant_id, r.name, o.status, o.total_kopecks, o.delivery_address, o.comment, o.created_at, o.updated_at
	FROM orders o
	JOIN restaurants r ON r.id = o.restaurant_id
	WHERE o.id = $1;
	`

	var order domain.OrderResponse

	err := r.pool.QueryRow(ctx, query, orderID).Scan(
		&order.ID,
		&order.UserID,
		&order.RestID,
		&order.RestName,
		&order.Status,
		&order.TotalKopecks,
		&order.DeliveryAddress,
		&order.Comment,
		&order.CreatedAt,
		&order.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.OrderResponse{}, fmt.Errorf("order not found: %w", core_errors.ErrNotFound)
	}
	if err != nil {
		return domain.OrderResponse{}, fmt.Errorf("scan order: %w", err)
	}

	items, err := r.getItems(ctx, []uuid.UUID{order.ID})
	if err != nil {
		return domain.OrderResponse{}, fmt.Errorf("get order items: %w", err)
	}

	order.Items = items[order.ID]
	if order.Items == nil {
		order.Items = []domain.OrderItem{}
	}

	return order, nil
}

func (r *OrdersRepository) getItems(
	ctx context.Context,
	orderIDs []uuid.UUID,
) (map[uuid.UUID][]domain.OrderItem, error) {
	query := `
	SELECT order_id, id, dish_id, dish_name, quantity, price_at_order_time
	FROM order_items
	WHERE order_id = ANY($1)
	ORDER BY dish_name;
	`

	rows, err := r.pool.Query(ctx, query, orderIDs)
	if err != nil {
		return nil, fmt.Errorf("query order items: %w", err)
	}
	defer rows.Close()

	items := make(map[uuid.UUID][]domain.OrderItem)

	for rows.Next() {
		var (
			orderID uuid.UUID
			item    domain.OrderItem
		)

		if err := rows.Scan(
			&orderID,
			&item.ID,
			&item.DishID,
			&item.DishName,
			&item.Quantity,
			&item.PriceAtOrderTime,
		); err != nil {
			return nil, fmt.Errorf("scan order item: %w", err)
		}

		items[orderID] = append(items[orderID], item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("order items rows: %w", err)
	}

	return items, nil
}
