package notifications_postgres_repository

import core_postgres_pool "github.com/PassDelivery/pass-backend/internal/core/repository/postgres/pool"

type NotificationsRepository struct {
	pool core_postgres_pool.Pool
}

func NewNotificationsRepository(
	pool core_postgres_pool.Pool,
) *NotificationsRepository {
	return &NotificationsRepository{
		pool: pool,
	}
}
