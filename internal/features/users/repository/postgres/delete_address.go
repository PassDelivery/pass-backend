package users_postgres_repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (r *UsersRepository) DeleteAddress(
	ctx context.Context,
	addrID uuid.UUID,
	user_id uuid.UUID,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	DELETE FROM addresses
	WHERE id = $1 AND user_id = $2
	`

	if _, err := r.pool.Exec(ctx, query, addrID, user_id); err != nil {
		return fmt.Errorf("failed exec to delet: %w", err)
	}

	return nil
}
