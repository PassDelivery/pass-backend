package users_postgres_repository

import (
	"time"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	"github.com/google/uuid"
)

type UserModel struct {
	ID        uuid.UUID
	Email     string
	Name      string
	Phone     *string
	Role      domain.Role
	CreatedAt time.Time
}

type UserModelWithHash struct {
	ID        uuid.UUID
	Email     string
	Hash      string
	Name      string
	Phone     *string
	Role      domain.Role
	CreatedAt time.Time
}
