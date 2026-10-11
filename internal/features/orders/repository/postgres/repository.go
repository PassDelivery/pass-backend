package orders_postgres_repository

import core_postgres_pool "github.com/PassDelivery/pass-backend/internal/core/repository/postgres/pool"

type OrdersRepository struct {
	pool core_postgres_pool.Pool
}

func NewOrdersRepository(
	pool core_postgres_pool.Pool,
) *OrdersRepository {
	return &OrdersRepository{
		pool: pool,
	}
}
