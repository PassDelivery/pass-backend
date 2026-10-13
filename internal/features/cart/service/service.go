package cart_service

import (
	"context"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	"github.com/google/uuid"
)

type CartService struct {
	cartRepository CartRepository
}

type CartRepository interface {
	GetCart(
		ctx context.Context,
		userID uuid.UUID,
	) (domain.Cart, error)
	AddItem(
		ctx context.Context,
		userID uuid.UUID,
		item domain.CartItemRequest,
	) error
	PatchItem(
		ctx context.Context,
		userID uuid.UUID,
		itemID uuid.UUID,
		quantity int,
	) error
	DeleteItem(
		ctx context.Context,
		userID uuid.UUID,
		itemID uuid.UUID,
	) error
	ClearCart(
		ctx context.Context,
		userID uuid.UUID,
	) error
}

func NewCartService(cr CartRepository) *CartService {
	return &CartService{
		cartRepository: cr,
	}
}
