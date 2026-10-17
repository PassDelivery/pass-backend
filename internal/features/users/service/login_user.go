package users_service

import (
	"context"
	"fmt"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	core_errors "github.com/PassDelivery/pass-backend/internal/core/errors"
	core_jwt "github.com/PassDelivery/pass-backend/internal/core/jwt"
)

func (s *UsersService) LoginUser(
	ctx context.Context,
	email string,
	password string,
) (domain.UserResponse, error) {
	if err := domain.ValidateEmail(email); err != nil {
		return domain.UserResponse{}, fmt.Errorf("validate email fo login: %w", err)
	}

	email, err := s.usersRepository.GetByEmail(ctx, email)
	if err != nil {
		return domain.UserResponse{}, fmt.Errorf("get email from repository: %w", err)
	}
	if email == "" {
		return domain.UserResponse{}, fmt.Errorf("invalid email or password: %w", core_errors.ErrUnauthtorize)
	}

	user, hash, err := s.usersRepository.GetUserByEmail(ctx, email)
	if err != nil {
		return domain.UserResponse{}, fmt.Errorf("get user by email: %w", err)
	}

	if ok := core_jwt.CheckPassword(password, hash); !ok {
		return domain.UserResponse{}, fmt.Errorf("invalid email or password: %w", core_errors.ErrUnauthtorize)
	}

	token, err := s.tokenService.Generate(user.ID, user.Role)
	if err != nil {
		return domain.UserResponse{}, fmt.Errorf("generate token: %w", err)
	}

	user.Token = token
	return user, nil
}
