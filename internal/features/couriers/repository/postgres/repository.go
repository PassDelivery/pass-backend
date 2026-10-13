package couriers_postgres_repository

import core_postgres_pool "github.com/PassDelivery/pass-backend/internal/core/repository/postgres/pool"

type CouriersRepository struct {
	pool core_postgres_pool.Pool
}

func NewCouriersRepository(
	pool core_postgres_pool.Pool,
) *CouriersRepository {
	return &CouriersRepository{
		pool: pool,
	}
}
