package domain

import (
	"fmt"
	"strings"
	"time"

	core_errors "github.com/PassDelivery/pass-backend/internal/core/errors"
	"github.com/google/uuid"
)

type Role string

const (
	RoleCustomer        Role = "customer"
	RoleRestaurantOwner Role = "restaurant_owner"
	RoleCourier         Role = "courier"
	RoleAdmin           Role = "admin"
	DefaultCost         int  = 10
)

func (r Role) IsValid() bool {
	switch r {
	case RoleCustomer, RoleRestaurantOwner, RoleCourier, RoleAdmin:
		return true
	}
	return false
}

type UserRequest struct {
	Email     string
	Password  string
	Name      string
	Phone     *string
	Role      Role
}

type UserResponse struct {
	ID        uuid.UUID
	Token     string
	Email     string
	Name      string
	Phone     *string
	Role      Role
	CreatedAt time.Time
}

type GetUserResponse struct {
	ID        uuid.UUID
	Email     string
	Name      string
	Phone     *string
	Role      Role
	CreatedAt time.Time
}

type UserWithHash struct {
	Email string
	Hash  string
	Name  string
	Phone *string
	Role  Role
}

type UserPatch struct {
	Email Nullable[string]
	Name  Nullable[string]
	Phone Nullable[string] 
}

func NewUserResponse(
	id uuid.UUID,
	token string,
	email string,
	name string,
	phone *string,
	role Role,
	createdAt time.Time,
) UserResponse {
	return UserResponse{
		ID:        id,
		Token:     token,
		Email:     email,
		Name:      name,
		Phone:     phone,
		Role:      role,
		CreatedAt: createdAt,
	}
}

func NewUserWithHash(
	email string,
	hash string,
	name string,
	phone *string,
	role Role,
) UserWithHash {
	return UserWithHash{
		Email: email,
		Hash:  hash,
		Name:  name,
		Phone: phone,
		Role:  role,
	}
}

func NewUser(
	email string,
	password string,
	name string,
	phone *string,
	role Role,
) UserRequest {
	return UserRequest{
		Email:     email,
		Password:  password,
		Name:      name,
		Phone:     phone,
		Role:      role,
	}
}



func NewUserUninitialized(
	email string,
	password string,
	name string,
	phone *string,
	role Role,
) UserRequest {
	return NewUser(
		email,
		password,
		name,
		phone,
		role,
	)
}

func (u *UserRequest) Validate() error {

	if err := ValidateEmail(u.Email); err != nil {
		return fmt.Errorf("validate email: %w", err)
	}

	if strings.TrimSpace(u.Name) == "" {
		return fmt.Errorf("name is required: %w", core_errors.ErrInvalidArgument)
	}

	if (*u.Phone)[0] != '+' || len(*u.Phone) < 12 {
		return fmt.Errorf("Number hasnt first + or len < 12 symbols: %w", core_errors.ErrInvalidArgument)
	}

	if !u.Role.IsValid() {
		return fmt.Errorf("Role is dont correct: %w", core_errors.ErrInvalidArgument)
	}

	return nil
}


func ValidateEmail(email string) error {
	email = strings.ToLower(email)
	if !strings.Contains(email, "@") {
		return fmt.Errorf("invalid email name: %w", core_errors.ErrInvalidArgument)
	}
	return nil
}