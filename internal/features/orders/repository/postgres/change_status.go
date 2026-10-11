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

func (r *OrdersRepository) ChangeOrderStatus(
	ctx context.Context,
	orderID uuid.UUID,
	from domain.OrderStatus,
	to domain.OrderStatus,
	changedBy uuid.UUID,
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

	queryUpdate := `
	UPDATE orders
	SET status = $1, updated_at = NOW()
	WHERE id = $2 AND status = $3
	RETURNING id, user_id, restaurant_id, status, total_kopecks, delivery_address, comment, created_at, updated_at;
	`

	queryHistory := `
	INSERT INTO order_status_history (order_id, status, changed_by, changed_at)
	VALUES ($1, $2, $3, NOW());
	`

	queryItems := `
	SELECT id, dish_id, dish_name, quantity, price_at_order_time
	FROM order_items
	WHERE order_id = $1;
	`

	var order domain.OrderResponse

	err = tx.QueryRow(ctx, queryUpdate, to, orderID, from).Scan(
		&order.ID,
		&order.UserID,
		&order.RestID,
		&order.Status,
		&order.TotalKopecks,
		&order.DeliveryAddress,
		&order.Comment,
		&order.CreatedAt,
		&order.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.OrderResponse{}, fmt.Errorf("order status already changed: %w", core_errors.ErrConflict)
	}
	if err != nil {
		return domain.OrderResponse{}, fmt.Errorf("failed scan: %w", err)
	}

	if _, err := tx.Exec(ctx, queryHistory, orderID, to, changedBy); err != nil {
		return domain.OrderResponse{}, fmt.Errorf("insert status history: %w", err)
	}

	rows, err := tx.Query(ctx, queryItems, orderID)
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

	if err := tx.Commit(ctx); err != nil {
		return domain.OrderResponse{}, fmt.Errorf("faild commit to database: %w", err)
	}

	return order, nil
}
