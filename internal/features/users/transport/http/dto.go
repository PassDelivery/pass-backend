package users_transport_http

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	core_http_types "github.com/PassDelivery/pass-backend/internal/core/transport/http/types"
	"github.com/google/uuid"
)

type LoginRequestDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type UserRequestDTO struct {
	Email    string  `json:"email" validate:"required,email"`
	Password string  `json:"password" validate:"required,min=6"`
	Name     string  `json:"name" validate:"required,min=2"`
	Phone    *string `json:"phone" validate:"omitempty,e164,min=10,max=15"`
	Role     string  `json:"role" validate:"required,oneof=customer restaurant_owner courier admin"`
}

type UserPatchDTO struct {
	Email core_http_types.Nullable[string] `json:"email"`
	Name  core_http_types.Nullable[string] `json:"name"`
	Phone core_http_types.Nullable[string] `json:"phone"`
}

type AddressPatchDTO struct {
	City      core_http_types.Nullable[string] `json:"city"`
	Street    core_http_types.Nullable[string] `json:"street"`
	House     core_http_types.Nullable[string] `json:"house"`
	Apartment core_http_types.Nullable[string] `json:"apartment"`
	Comment   core_http_types.Nullable[string] `json:"comment"`
	IsDefault core_http_types.Nullable[bool] `json:"is_default"`
}

type UserResponseDTO struct {
	Token     string    `json:"token"`
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Phone     *string   `json:"phone"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

type GetUserResponseDTO struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Phone     *string   `json:"phone"`
	Role      string    `json:"role"`
	AvatarURL *string   `json:"avatar_url"`
	CreatedAt time.Time `json:"created_at"`
}

type AddressesResponseDTO struct {
	City      string    `json:"city"`
	Street    string    `json:"street"`
	House     string    `json:"house"`
	Apartment *string   `json:"apartment"`
	Comment   *string   `json:"comment"`
	IsDefault bool      `json:"is_default"`
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

type AddressRequestDTO struct {
	City      string  `json:"city" validate:"required,min=1"`
	Street    string  `json:"street" validate:"required,min=1"`
	House     string  `json:"house" validate:"required,min=1"`
	Apartment *string `json:"apartment" validate:"omitempty,min=1"`
	Comment   *string `json:"comment" validate:"omitempty,min=1,max=50"`
	IsDefault bool    `json:"is_default"`
}

func addressRequestDomainFromDTO(req AddressRequestDTO) domain.AddressRequest {
	return domain.AddressRequest{
		City:      req.City,
		Street:    req.Street,
		House:     req.House,
		Apartment: req.Apartment,
		Comment:   req.Comment,
		IsDefault: req.IsDefault,
	}
}

func addressResponseDTOFromDomain(res domain.AddressesResponse) AddressesResponseDTO {
	return AddressesResponseDTO{
		City:      res.City,
		Street:    res.Street,
		House:     res.House,
		Apartment: res.Apartment,
		Comment:   res.Comment,
		IsDefault: res.IsDefault,
		ID:        res.ID,
		UserID:    res.UserID,
		CreatedAt: res.CreatedAt,
	}
}

func NewAddressesResponseDTO(addr []domain.AddressesResponse) []AddressesResponseDTO {
	res := make([]AddressesResponseDTO, len(addr))
	for k, v := range addr {
		res[k] = AddressesResponseDTO{
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

func domainFromDTO(dto UserRequestDTO) domain.UserRequest {
	return domain.NewUserUninitialized(
		dto.Email,
		dto.Password,
		dto.Name,
		dto.Phone,
		domain.Role(dto.Role),
	)
}

func dtoFromDomain(user domain.UserResponse) UserResponseDTO {
	return UserResponseDTO{
		Token:     user.Token,
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Phone:     user.Phone,
		Role:      string(user.Role),
		CreatedAt: user.CreatedAt,
	}
}

func (d *UserPatchDTO) Validate() error {
	if d.Email.Set {
		if d.Email.Value == nil {
			return fmt.Errorf("Email cant be null")
		}
		if ok := strings.Contains(*d.Email.Value, "@"); !ok {
			return fmt.Errorf("invalid symbol in Email")
		}
		if len(*d.Email.Value) < 5 {
			return fmt.Errorf("Email symbols count cant be < 5")
		}
	}

	if d.Name.Set {
		if d.Name.Value == nil {
			return fmt.Errorf("Name cant be null")
		}
		if len(*d.Name.Value) < 2 {
			return fmt.Errorf("Name symbols count cant be < 2")
		}
	}

	if d.Phone.Set {
		if d.Phone.Value == nil || *d.Phone.Value == "" {
			return nil
		}
		if ok := strings.Contains(*d.Phone.Value, "+"); !ok {
			return fmt.Errorf("invalid symbol in Phone")
		}
		if len(*d.Phone.Value) < 10 {
			return fmt.Errorf("Phone symbols count cant be < 10")
		}
		if len(*d.Phone.Value) > 15 {
			return fmt.Errorf("Phone symbols count cant be > 15")
		}
	}

	return nil
}

func (d *AddressPatchDTO) Validate() error {
	if d.City.Set {
		if d.City.Value == nil {
			return fmt.Errorf("city cant be null")
		}
		if len(*d.City.Value) <= 2 {
			return fmt.Errorf("len city cant be < 3")
		}
	}

	if d.Street.Set {
		if d.Street.Value == nil {
			return fmt.Errorf("street cant be null")
		}
		if len(*d.Street.Value) <= 2 {
			return fmt.Errorf("len street cant be < 3")
		}
	}

	if d.House.Set {
		if d.House.Value == nil {
			return fmt.Errorf("house cant be null")
		}
		_, err := strconv.Atoi(*d.House.Value)
		if err != nil {
			return fmt.Errorf("house number is word")
		}
	}

	if d.Apartment.Set {
		if d.Apartment.Value == nil {
			return fmt.Errorf("apartment cant be null")
		}
		if *d.Apartment.Value != "" {
			_, err := strconv.Atoi(*d.Apartment.Value)
			if err != nil {
				return fmt.Errorf("apartment number is word")
			}
		}
	}

	if d.Comment.Set {
		if d.Comment.Value == nil {
			return fmt.Errorf("comment cant be null")
		}
		if *d.Comment.Value != "" {
			if len(*d.Comment.Value) < 5 {
				return fmt.Errorf("len comment is < 5")
			}
			if len(*d.Comment.Value) > 100 {
				return fmt.Errorf("len comment is > 100")
			}
		}
	}

	if d.IsDefault.Set {
		if d.IsDefault.Value == nil {
			return fmt.Errorf("is_default cant be null")
		}
	}
	return nil
}