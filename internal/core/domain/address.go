package domain

import (
	"fmt"
	"time"

	core_errors "github.com/PassDelivery/pass-backend/internal/core/errors"
	"github.com/google/uuid"
)

type AddressesResponse struct {
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

type AddressRequest struct {
	City      string
	Street    string
	House     string
	Apartment *string
	Comment   *string
	IsDefault bool
}

type AddressPatch struct {
	City      Nullable[string]
	Street    Nullable[string]
	House     Nullable[string]
	Apartment Nullable[string]
	Comment   Nullable[string]
	IsDefault Nullable[bool]
}

func (p *AddressPatch) Validate() error {
	if p.City.Set && p.City.Value == nil {
		return fmt.Errorf("Email dont be nil: %w", core_errors.ErrInvalidArgument)
	}
	if p.Street.Set && p.Street.Value == nil {
		return fmt.Errorf("Name dont be nil: %w", core_errors.ErrInvalidArgument)
	}
	if p.House.Set && p.House.Value == nil {
		return fmt.Errorf("Name dont be nil: %w", core_errors.ErrInvalidArgument)
	}
	if p.IsDefault.Set && p.IsDefault.Value == nil {
		return fmt.Errorf("Name dont be nil: %w", core_errors.ErrInvalidArgument)
	}
	return nil
}

func (u *AddressRequest) ApplyPatch(patch AddressPatch) error {
	if err := patch.Validate(); err != nil {
		return fmt.Errorf("validate address patch: %w", err)
	}

	tmp := *u

	if patch.City.Set {
		tmp.City = *patch.City.Value
	}
	if patch.Street.Set {
		tmp.Street = *patch.Street.Value
	}
	if patch.House.Set {
		tmp.House = *patch.House.Value
	}
	if patch.IsDefault.Set {
		tmp.IsDefault = *patch.IsDefault.Value
	}
	if patch.Apartment.Set {
		tmp.Apartment = patch.Apartment.Value
	}
	if patch.Comment.Set {
		tmp.Comment = patch.Comment.Value
	}

	*u = tmp

	return nil
}
