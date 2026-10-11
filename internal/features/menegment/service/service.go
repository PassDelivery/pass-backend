package menegment_service

import (
	"context"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	"github.com/google/uuid"
)

type MenegmentService struct {
	menegmentRepository MenegmentRepository
}

type MenegmentRepository interface {
	RegisterRestaurant(
		ctx context.Context,
		restaurant domain.RestaurantRequest,
		id uuid.UUID,
	) (domain.RestaurantResponse, error)
	GetRestaurant(
	ctx context.Context,
	restID uuid.UUID,
	) (domain.RestaurantRequest, error)
	PatchRestaurant(
		ctx context.Context,
		rest domain.RestaurantRequest,
		restID uuid.UUID,
		ownerID uuid.UUID,
	) (domain.RestaurantResponse, error)
	AddCategory(
		ctx context.Context,
		category domain.CategoryRequest,
		restID uuid.UUID,
	) (domain.CategoryResponse, error)
	IsRestaurantOwner(
		ctx context.Context,
		restID uuid.UUID,
		ownerID uuid.UUID,
	) (bool, error)
	IsCategoryInRestaurant(
		ctx context.Context,
		catID uuid.UUID,
		restID uuid.UUID,
	) (bool, error)
	AddDish(
		ctx context.Context,
		dish domain.DishRequest,
		restID uuid.UUID,
	) (domain.DishResponse, error)
	GetDish(
		ctx context.Context,
		dishID uuid.UUID,
	) (domain.DishResponse, error)
	PatchDish(
		ctx context.Context,
		dish domain.DishRequest,
		dishID uuid.UUID,
	) (domain.DishResponse, error)
	IsDishInActiveOrders(
		ctx context.Context,
		dishID uuid.UUID,
	) (bool, error)
	HideDish(
		ctx context.Context,
		dishID uuid.UUID,
	) error
}

func NewMenegmentService(mr MenegmentRepository) *MenegmentService {
	return &MenegmentService{
		menegmentRepository: mr,
	}
}
