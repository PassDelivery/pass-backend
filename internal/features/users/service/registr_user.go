package users_service

import (
	"context"
	"fmt"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	core_errors "github.com/PassDelivery/pass-backend/internal/core/errors"
	"golang.org/x/crypto/bcrypt"
)

func (s *UsersService) RegisterUser(
	ctx context.Context,
	user domain.UserRequest,
) (domain.UserResponse, error) {
	if err := user.Validate(); err != nil {
		return domain.UserResponse{}, fmt.Errorf("validate user domain: %w", err)
	}

	exist, err := s.usersRepository.GetByEmail(ctx, user.Email)

	if err != nil {
		return domain.UserResponse{}, fmt.Errorf("failed get email: %w", err)
	}

	if exist != "" {
		return domain.UserResponse{}, fmt.Errorf("email already taken: %w", core_errors.ErrConflict)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), domain.DefaultCost)
	if err != nil {
		return domain.UserResponse{}, fmt.Errorf("generate hash from password: %w", err)
	}

	userWithHash := domain.NewUserWithHash(
		user.Email,
		string(hash),
		user.Name,
		user.Phone,
		user.Role,
	)

	userRes, err := s.usersRepository.CreateUser(ctx, userWithHash)
	if err != nil {
		return domain.UserResponse{}, fmt.Errorf("create user: %w", err)
	}

	token, err := s.tokenService.Generate(userRes.ID, userRes.Role)
	if err != nil {
		return domain.UserResponse{}, fmt.Errorf("generate token: %w", err)
	}

	userResponse := domain.NewUserResponse(
		userRes.ID,
		token,
		userRes.Email,
		userRes.Name,
		userRes.Phone,
		userRes.Role,
		userRes.CreatedAt,
	)

	return userResponse, nil
}
