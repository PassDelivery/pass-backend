package users_postgres_repository

import (
	"time"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	"github.com/google/uuid"
)

type UserModel struct {
	ID        uuid.UUID
	Email     string
	Name      string
	Phone     *string
	Role      domain.Role
	CreatedAt time.Time
}

type UserModelWithHash struct {
	ID        uuid.UUID
	Email     string
	Hash      string
	Name      string
	Phone     *string
	Role      domain.Role
	CreatedAt time.Time
}

type AddressesResponseModel struct {
	City      string
	Street    string
	House     string
	Apartment *string
	Comment   *string
	IsDefault bool
	ID        uuid.UUID
	UserID    uuid.UUID
	CreatedAt time.Time
}

func addressDomainFromModel(model AddressesResponseModel) domain.AddressesResponse {
	return domain.AddressesResponse{
		City:      model.City,
		Street:    model.Street,
		House:     model.House,
		Apartment: model.Apartment,
		Comment:   model.Comment,
		IsDefault: model.IsDefault,
		ID:        model.ID,
		UserID:    model.UserID,
		CreatedAt: model.CreatedAt,
	}
}

func addressDomainsFromModels(model []AddressesResponseModel) []domain.AddressesResponse {
	res := make([]domain.AddressesResponse, len(model))

	for k, v := range model {
			res[k] = domain.AddressesResponse{
				City:      v.City,
				Street:    v.Street,
				House:     v.House,
				Apartment: v.Apartment,
				Comment:   v.Comment,
				IsDefault: v.IsDefault,
				ID:        v.ID,
				UserID:    v.UserID,
				CreatedAt: v.CreatedAt,
			}
	}

	return res
}
