package orders_postgres_repository

import (
	"context"
	"fmt"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	"github.com/google/uuid"
)

func (r *OrdersRepository) GetUserOrders(
	ctx context.Context,
	userID uuid.UUID,
) ([]domain.OrderResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT o.id, o.user_id, o.restaurant_id, r.name, o.status, o.total_kopecks, o.delivery_address, o.comment, o.created_at, o.updated_at
	FROM orders o
	JOIN restaurants r ON r.id = o.restaurant_id
	WHERE o.user_id = $1
	ORDER BY o.created_at DESC;
	`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("query orders: %w", err)
	}

	orders := []domain.OrderResponse{}
	ids := []uuid.UUID{}

	for rows.Next() {
		var order domain.OrderResponse

		if err := rows.Scan(
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
		); err != nil {
			rows.Close()
			return nil, fmt.Errorf("scan order: %w", err)
		}

		orders = append(orders, order)
		ids = append(ids, order.ID)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("orders rows: %w", err)
	}

	if len(ids) == 0 {
		return orders, nil
	}

	queryItems := `
	SELECT order_id, id, dish_id, dish_name, quantity, price_at_order_time
	FROM order_items
	WHERE order_id = ANY($1)
	ORDER BY dish_name;
	`

	itemRows, err := r.pool.Query(ctx, queryItems, ids)
	if err != nil {
		return nil, fmt.Errorf("query order items: %w", err)
	}
	defer itemRows.Close()

	items := make(map[uuid.UUID][]domain.OrderItem)

	for itemRows.Next() {
		var (
			orderID uuid.UUID
			item    domain.OrderItem
		)

		if err := itemRows.Scan(
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
	if err := itemRows.Err(); err != nil {
		return nil, fmt.Errorf("order items rows: %w", err)
	}

	for k := range orders {
		orders[k].Items = items[orders[k].ID]
		if orders[k].Items == nil {
			orders[k].Items = []domain.OrderItem{}
		}
	}

	return orders, nil
}
