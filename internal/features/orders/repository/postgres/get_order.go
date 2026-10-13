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

	queryItems := `
	SELECT id, dish_id, dish_name, quantity, price_at_order_time
	FROM order_items
	WHERE order_id = $1
	ORDER BY dish_name;
	`

	rows, err := r.pool.Query(ctx, queryItems, order.ID)
	if err != nil {
		return domain.OrderResponse{}, fmt.Errorf("query order items: %w", err)
	}
	defer rows.Close()

	order.Items = []domain.OrderItem{}

	for rows.Next() {
		var item domain.OrderItem

		if err := rows.Scan(
			&item.ID,
			&item.DishID,
			&item.DishName,
			&item.Quantity,
			&item.PriceAtOrderTime,
		); err != nil {
			return domain.OrderResponse{}, fmt.Errorf("scan order item: %w", err)
		}

		order.Items = append(order.Items, item)
	}
	if err := rows.Err(); err != nil {
		return domain.OrderResponse{}, fmt.Errorf("order items rows: %w", err)
	}

	return order, nil
}
