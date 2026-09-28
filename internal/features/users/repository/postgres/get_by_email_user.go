package users_postgres_repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func (r *UsersRepository) GetByEmail(
	ctx context.Context,
	email string,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT email
	FROM users
	WHERE LOWER(email) = LOWER($1);
	`
	var mail string

	if err := r.pool.QueryRow(ctx, query, email).Scan(&mail); err == pgx.ErrNoRows {
		return "", nil
	}else if err != nil {
		return "", fmt.Errorf("failed scan from database: %w", err)
	}

	return mail, nil
}
