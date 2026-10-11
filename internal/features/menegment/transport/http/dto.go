package menegment_transport_http

import (
	"fmt"
	"time"

	core_errors "github.com/PassDelivery/pass-backend/internal/core/errors"
	core_http_types "github.com/PassDelivery/pass-backend/internal/core/transport/http/types"
	"github.com/google/uuid"
)

type RegisterRestaurantRequestDTO struct {
	Name        string  `json:"name" validate:"required,min=2"`
	Description *string `json:"description" validate:"omitempty,min=5,max=100"`
	Address     string  `json:"address" validate:"required,min=5"`
	Phone       *string `json:"phone" validate:"omitempty,e164,min=10,max=15"`
	IsOpen      bool    `json:"is_open" validate:"required"`
}

type RestaurantResponseDTO struct {
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	Address     string    `json:"address"`
	Phone       *string   `json:"phone"`
	IsOpen      bool      `json:"is_open"`
	ID          uuid.UUID `json:"id"`
	OwnerID     uuid.UUID `json:"owner_id"`
	CreatedAt   time.Time `json:"created_at"`
}

type RestaurantPatchDTO struct {
	Name        core_http_types.Nullable[string] `json:"name"`
	Description core_http_types.Nullable[string] `json:"description"`
	Address     core_http_types.Nullable[string] `json:"address"`
	Phone       core_http_types.Nullable[string] `json:"phone"`
	IsOpen      core_http_types.Nullable[bool]   `json:"is_open"`
}

func (p *RestaurantPatchDTO) Validate() error {
	if p.Name.Set {
		if *p.Name.Value == "" || p.Name.Value == nil {
			return fmt.Errorf("%w", core_errors.ErrInvalidArgument)
		}
		if len([]rune(*p.Name.Value)) < 2 {
			return fmt.Errorf("%w", core_errors.ErrInvalidArgument)
		}
	}
	if p.Description.Set {
		if *p.Description.Value != "" {
			if len([]rune(*p.Description.Value)) < 5 || len([]rune(*p.Description.Value)) > 100 {
				return fmt.Errorf("%w", core_errors.ErrInvalidArgument)
			}
		}
	}
	if p.Address.Set {
		if *p.Address.Value == "" || p.Address.Value == nil {
			return fmt.Errorf("%w", core_errors.ErrInvalidArgument)
		}
		if len([]rune(*p.Address.Value)) < 5 {
			return fmt.Errorf("%w", core_errors.ErrInvalidArgument)
		}
	}
	if p.Phone.Set {
		if *p.Phone.Value != "" {
			if len([]rune(*p.Phone.Value)) < 10 || len([]rune(*p.Phone.Value)) > 15 {
				return fmt.Errorf("%w", core_errors.ErrInvalidArgument)
			}
			if []rune(*p.Phone.Value)[0] != '+' || []rune(*p.Phone.Value)[1] != '7' {
				return fmt.Errorf("%w", core_errors.ErrInvalidArgument)
			}
		}
	}
	if p.IsOpen.Set {
		if p.IsOpen.Value == nil {
			return fmt.Errorf("%w", core_errors.ErrInvalidArgument)
		}
	}

	return nil
}

type CategoryRequestDTO struct {
	Name      string `json:"name" validate:"required,min=3,max=35"`
	SortOrder int    `json:"sort_order" validate:"required"`
}

type CategoryResponseDTO struct {
	ID        uuid.UUID `json:"id"`
	RestID    uuid.UUID `json:"restaurant_id"`
	Name      string    `json:"name"`
	SortOrder string    `json:"sort_order"`
}

type DishRequestDTO struct {
	CategoryID         string  `json:"category_id" validate:"required,uuid"`
	Name               string  `json:"name" validate:"required,min=2,max=50"`
	Description        *string `json:"description" validate:"omitempty,min=5,max=200"`
	PriceKopecks       int64   `json:"price_kopecks" validate:"required,gt=0"`
	CookingTimeMinutes *int    `json:"cooking_time_minutes" validate:"omitempty,gt=0"`
	IsAvailable        bool    `json:"is_available"`
}

type DishResponseDTO struct {
	ID                 uuid.UUID `json:"id"`
	RestID             uuid.UUID `json:"restaurant_id"`
	CategoryID         uuid.UUID `json:"category_id"`
	Name               string    `json:"name"`
	Description        *string   `json:"description"`
	PriceKopecks       int64     `json:"price_kopecks"`
	CookingTimeMinutes *int      `json:"cooking_time_minutes"`
	IsAvailable        bool      `json:"is_available"`
	CreatedAt          time.Time `json:"created_at"`
}

type DishPatchDTO struct {
	CategoryID         core_http_types.Nullable[string] `json:"category_id"`
	Name               core_http_types.Nullable[string] `json:"name"`
	Description        core_http_types.Nullable[string] `json:"description"`
	PriceKopecks       core_http_types.Nullable[int64]  `json:"price_kopecks"`
	CookingTimeMinutes core_http_types.Nullable[int]    `json:"cooking_time_minutes"`
	IsAvailable        core_http_types.Nullable[bool]   `json:"is_available"`
}

func (p *DishPatchDTO) Validate() error {
	if p.CategoryID.Set {
		if p.CategoryID.Value == nil {
			return fmt.Errorf("category_id cant be null: %w", core_errors.ErrInvalidArgument)
		}
		if _, err := uuid.Parse(*p.CategoryID.Value); err != nil {
			return fmt.Errorf("category_id is not uuid: %w", core_errors.ErrInvalidArgument)
		}
	}
	if p.Name.Set {
		if p.Name.Value == nil {
			return fmt.Errorf("name cant be null: %w", core_errors.ErrInvalidArgument)
		}
		if len([]rune(*p.Name.Value)) < 2 || len([]rune(*p.Name.Value)) > 50 {
			return fmt.Errorf("name len must be 2..50: %w", core_errors.ErrInvalidArgument)
		}
	}
	if p.Description.Set && p.Description.Value != nil {
		if len([]rune(*p.Description.Value)) < 5 || len([]rune(*p.Description.Value)) > 200 {
			return fmt.Errorf("description len must be 5..200: %w", core_errors.ErrInvalidArgument)
		}
	}
	if p.PriceKopecks.Set {
		if p.PriceKopecks.Value == nil {
			return fmt.Errorf("price_kopecks cant be null: %w", core_errors.ErrInvalidArgument)
		}
		if *p.PriceKopecks.Value <= 0 {
			return fmt.Errorf("price_kopecks must be > 0: %w", core_errors.ErrInvalidArgument)
		}
	}
	if p.CookingTimeMinutes.Set && p.CookingTimeMinutes.Value != nil {
		if *p.CookingTimeMinutes.Value <= 0 {
			return fmt.Errorf("cooking_time_minutes must be > 0: %w", core_errors.ErrInvalidArgument)
		}
	}
	if p.IsAvailable.Set {
		if p.IsAvailable.Value == nil {
			return fmt.Errorf("is_available cant be null: %w", core_errors.ErrInvalidArgument)
		}
	}

	return nil
}
