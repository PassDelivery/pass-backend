package couriers_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	core_errors "github.com/PassDelivery/pass-backend/internal/core/errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *CouriersRepository) GetDelivery(
	ctx context.Context,
	courierID uuid.UUID,
	deliveryID uuid.UUID,
) (domain.Delivery, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	queryDelivery := `
	SELECT d.id, d.status, d.offered_at, d.responded_at, d.delivered_at,
		o.id, o.status, r.name, r.address, o.delivery_address, o.comment, u.name, u.phone, o.total_kopecks
	FROM deliveries d
	JOIN orders o ON o.id = d.order_id
	JOIN restaurants r ON r.id = o.restaurant_id
	JOIN users u ON u.id = o.user_id
	WHERE d.courier_id = $1 AND d.id = $2;
	`

	queryItems := `
	SELECT id, dish_id, dish_name, quantity, price_at_order_time
	FROM order_items
	WHERE order_id = $1
	ORDER BY dish_name;
	`

	var delivery domain.Delivery

	err := r.pool.QueryRow(ctx, queryDelivery, courierID, deliveryID).Scan(
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
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Delivery{}, fmt.Errorf("delivery not found: %w", core_errors.ErrNotFound)
	}
	if err != nil {
		return domain.Delivery{}, fmt.Errorf("scan delivery: %w", err)
	}

	rows, err := r.pool.Query(ctx, queryItems, delivery.OrderID)
	if err != nil {
		return domain.Delivery{}, fmt.Errorf("query order items: %w", err)
	}
	defer rows.Close()

	delivery.Items = []domain.OrderItem{}

	for rows.Next() {
		var item domain.OrderItem

		if err := rows.Scan(
			&item.ID,
			&item.DishID,
			&item.DishName,
			&item.Quantity,
			&item.PriceAtOrderTime,
		); err != nil {
			return domain.Delivery{}, fmt.Errorf("scan order item: %w", err)
		}

		delivery.Items = append(delivery.Items, item)
	}
	if err := rows.Err(); err != nil {
		return domain.Delivery{}, fmt.Errorf("order items rows: %w", err)
	}

	return delivery, nil
}
