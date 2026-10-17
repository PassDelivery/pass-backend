package domain

import (
	"fmt"
	"time"

	core_errors "github.com/PassDelivery/pass-backend/internal/core/errors"
	"github.com/google/uuid"
)

type RestaurantRequest struct {
	Name        string
	Description *string
	Address     string
	Phone       *string
	IsOpen      bool
}

type RestaurantResponse struct {
	Name        string
	Description *string
	Address     string
	Phone       *string
	IsOpen      bool
	PhotoURL    *string
	ID          uuid.UUID
	OwnerID     uuid.UUID
	CreatedAt   time.Time
}

type RestaurantPatch struct {
	Name        Nullable[string] 
	Description Nullable[string]
	Address     Nullable[string]
	Phone       Nullable[string]
	IsOpen      Nullable[bool]
}

type CategoryRequest struct {
	Name      string 
	SortOrder int    
}

type CategoryResponse struct {
	ID        uuid.UUID 
	RestID    uuid.UUID
	Name      string 
	SortOrder string 
}

type DishRequest struct {
	CategoryID         uuid.UUID
	Name               string
	Description        *string
	PriceKopecks       int64
	CookingTimeMinutes *int
	IsAvailable        bool
}

type DishResponse struct {
	ID                 uuid.UUID
	RestID             uuid.UUID
	CategoryID         uuid.UUID
	Name               string
	Description        *string
	PriceKopecks       int64
	CookingTimeMinutes *int
	IsAvailable        bool
	PhotoURL           *string
	CreatedAt          time.Time
}

type DishPatch struct {
	CategoryID         Nullable[uuid.UUID]
	Name               Nullable[string]
	Description        Nullable[string]
	PriceKopecks       Nullable[int64]
	CookingTimeMinutes Nullable[int]
	IsAvailable        Nullable[bool]
}



func (r *RestaurantRequest) ApplyPatch(patch RestaurantPatch){
	tmp := *r
	
	if patch.Name.Set {
		tmp.Name = *patch.Name.Value
	}
	if patch.Description.Set {
		tmp.Description = patch.Description.Value
	}
	if patch.Address.Set {
		tmp.Address = *patch.Address.Value
	}
	if patch.Phone.Set {
		tmp.Phone = patch.Phone.Value
	}
	if patch.IsOpen.Set {
		tmp.IsOpen = *patch.IsOpen.Value
	}

	*r = tmp
}

func (p *DishPatch) Validate() error {
	if p.CategoryID.Set && p.CategoryID.Value == nil {
		return fmt.Errorf("category_id cant be null: %w", core_errors.ErrInvalidArgument)
	}
	if p.Name.Set && p.Name.Value == nil {
		return fmt.Errorf("name cant be null: %w", core_errors.ErrInvalidArgument)
	}
	if p.PriceKopecks.Set && p.PriceKopecks.Value == nil {
		return fmt.Errorf("price_kopecks cant be null: %w", core_errors.ErrInvalidArgument)
	}
	if p.IsAvailable.Set && p.IsAvailable.Value == nil {
		return fmt.Errorf("is_available cant be null: %w", core_errors.ErrInvalidArgument)
	}
	return nil
}

func (d *DishRequest) ApplyPatch(patch DishPatch) error {
	if err := patch.Validate(); err != nil {
		return fmt.Errorf("validate dish patch: %w", err)
	}

	tmp := *d

	if patch.CategoryID.Set {
		tmp.CategoryID = *patch.CategoryID.Value
	}
	if patch.Name.Set {
		tmp.Name = *patch.Name.Value
	}
	if patch.Description.Set {
		tmp.Description = patch.Description.Value
	}
	if patch.PriceKopecks.Set {
		tmp.PriceKopecks = *patch.PriceKopecks.Value
	}
	if patch.CookingTimeMinutes.Set {
		tmp.CookingTimeMinutes = patch.CookingTimeMinutes.Value
	}
	if patch.IsAvailable.Set {
		tmp.IsAvailable = *patch.IsAvailable.Value
	}

	*d = tmp

	return nil
}
