package menegment_postgres_repository

import core_postgres_pool "github.com/PassDelivery/pass-backend/internal/core/repository/postgres/pool"

type MenegmentRepository struct {
	pool core_postgres_pool.Pool
}


func NewMenegmentRepository(
	pool core_postgres_pool.Pool,
) *MenegmentRepository {
	return &MenegmentRepository{
		pool: pool,
	}
}