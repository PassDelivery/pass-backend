package users_service

import (
	"context"
	"fmt"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	core_files "github.com/PassDelivery/pass-backend/internal/core/files"
	core_jwt "github.com/PassDelivery/pass-backend/internal/core/jwt"
)

func (s *UsersService) UploadAvatar(
	ctx context.Context,
	data []byte,
) (domain.GetUserResponse, error) {
	userID, err := core_jwt.GetIDFromContext(ctx)
	if err != nil {
		return domain.GetUserResponse{}, fmt.Errorf("get id from context: %w", err)
	}

	url, err := core_files.SavePhoto(core_files.PhotoKindAvatars, data)
	if err != nil {
		return domain.GetUserResponse{}, fmt.Errorf("save avatar: %w", err)
	}

	user, oldURL, err := s.usersRepository.SetAvatar(ctx, userID, url)
	if err != nil {
		core_files.RemovePhoto(&url)
		return domain.GetUserResponse{}, fmt.Errorf("set avatar in service: %w", err)
	}

	core_files.RemovePhoto(oldURL)

	return user, nil
}
