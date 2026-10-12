package orders_postgres_repository

import (
	"context"
	"fmt"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	"github.com/google/uuid"
)

func (r *OrdersRepository) IsRestaurantOwner(
	ctx context.Context,
	restID uuid.UUID,
	ownerID uuid.UUID,
) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT EXISTS(
		SELECT 1 FROM restaurants
		WHERE id = $1 AND owner_id = $2
	);
	`

	var ok bool

	if err := r.pool.QueryRow(ctx, query, restID, ownerID).Scan(&ok); err != nil {
		return false, fmt.Errorf("failed scan: %w", err)
	}

	return ok, nil
}

func (r *OrdersRepository) GetRestaurantOrders(
	ctx context.Context,
	restID uuid.UUID,
) ([]domain.OrderResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT o.id, o.user_id, u.name, u.phone, o.restaurant_id, r.name, o.status, o.total_kopecks, o.delivery_address, o.comment, o.created_at, o.updated_at
	FROM orders o
	JOIN restaurants r ON r.id = o.restaurant_id
	JOIN users u ON u.id = o.user_id
	WHERE o.restaurant_id = $1
	ORDER BY o.created_at DESC;
	`

	rows, err := r.pool.Query(ctx, query, restID)
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
			&order.CustomerName,
			&order.CustomerPhone,
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

	items, err := r.getItems(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("get orders items: %w", err)
	}

	for k := range orders {
		orders[k].Items = items[orders[k].ID]
		if orders[k].Items == nil {
			orders[k].Items = []domain.OrderItem{}
		}
	}

	return orders, nil
}
