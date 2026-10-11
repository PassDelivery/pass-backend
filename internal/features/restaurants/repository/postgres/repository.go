package restaurants_postgres_repository

import core_postgres_pool "github.com/PassDelivery/pass-backend/internal/core/repository/postgres/pool"

type RestaurantsRepository struct {
	pool core_postgres_pool.Pool
}

func NewRestaurantsRepository(
	pool core_postgres_pool.Pool,
) *RestaurantsRepository {
	return &RestaurantsRepository{
		pool: pool,
	}
}
