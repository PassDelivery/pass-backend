package users_service

import (
	"context"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	core_jwt "github.com/PassDelivery/pass-backend/internal/core/jwt"
	"github.com/google/uuid"
)

type UsersService struct {
	usersRepository UsersRepository
	tokenService    core_jwt.TokenService
}

type UsersRepository interface {
	CreateUser(
		ctx context.Context,
		user domain.UserWithHash,
	) (domain.UserResponse, error)
	GetByEmail(
		ctx context.Context,
		email string,
	) (string, error)
	GetUserByEmail(
		ctx context.Context,
		email string,
	) (domain.UserResponse, string, error)
	GetUser(
		ctx context.Context,
		id uuid.UUID,
	) (domain.GetUserResponse, error)
}

func NewUsersService(
	usersRepository UsersRepository,
	tokenService core_jwt.TokenService,
) *UsersService {
	return &UsersService{
		usersRepository: usersRepository,
		tokenService:    tokenService,
	}
}
