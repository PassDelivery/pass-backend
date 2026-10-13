package couriers_postgres_repository

import (
	"context"
	"fmt"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	"github.com/google/uuid"
)

func (r *CouriersRepository) GetDeliveries(
	ctx context.Context,
	courierID uuid.UUID,
) ([]domain.Delivery, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	queryDeliveries := `
	SELECT d.id, d.status, d.offered_at, d.responded_at, d.delivered_at,
		o.id, o.status, r.name, r.address, o.delivery_address, o.comment, u.name, u.phone, o.total_kopecks
	FROM deliveries d
	JOIN orders o ON o.id = d.order_id
	JOIN restaurants r ON r.id = o.restaurant_id
	JOIN users u ON u.id = o.user_id
	WHERE d.courier_id = $1 AND d.status IN ('offered', 'accepted', 'picked_up')
	ORDER BY d.offered_at;
	`

	queryItems := `
	SELECT order_id, id, dish_id, dish_name, quantity, price_at_order_time
	FROM order_items
	WHERE order_id = ANY($1)
	ORDER BY dish_name;
	`

	rows, err := r.pool.Query(ctx, queryDeliveries, courierID)
	if err != nil {
		return nil, fmt.Errorf("query deliveries: %w", err)
	}

	deliveries := []domain.Delivery{}
	orderIDs := []uuid.UUID{}

	for rows.Next() {
		var delivery domain.Delivery

		if err := rows.Scan(
			&delivery.ID,
			&delivery.Status,
			&delivery.OfferedAt,
			&delivery.RespondedAt,
			&delivery.DeliveredAt,
			&delivery.OrderID,
			&delivery.OrderStatus,
			&delivery.RestName,
			&delivery.RestAddress,
			&delivery.DeliveryAddress,
			&delivery.Comment,
			&delivery.CustomerName,
			&delivery.CustomerPhone,
			&delivery.TotalKopecks,
		); err != nil {
			rows.Close()
			return nil, fmt.Errorf("scan delivery: %w", err)
		}

		deliveries = append(deliveries, delivery)
		orderIDs = append(orderIDs, delivery.OrderID)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("deliveries rows: %w", err)
	}

	if len(orderIDs) == 0 {
		return deliveries, nil
	}

	itemRows, err := r.pool.Query(ctx, queryItems, orderIDs)
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

	for k := range deliveries {
		deliveries[k].Items = items[deliveries[k].OrderID]
		if deliveries[k].Items == nil {
			deliveries[k].Items = []domain.OrderItem{}
		}
	}

	return deliveries, nil
}
