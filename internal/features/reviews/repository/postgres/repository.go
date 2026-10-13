package reviews_postgres_repository

import core_postgres_pool "github.com/PassDelivery/pass-backend/internal/core/repository/postgres/pool"

type ReviewsRepository struct {
	pool core_postgres_pool.Pool
}

func NewReviewsRepository(
	pool core_postgres_pool.Pool,
) *ReviewsRepository {
	return &ReviewsRepository{
		pool: pool,
	}
}
