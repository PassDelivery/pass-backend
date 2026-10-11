package users_transport_http

import (
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

type UserResponseDTO struct {
	Token     string    `json:"token"`
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Phone     *string   `json:"phone"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
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
