package repository

import (
	"context"

	"github.com/xenptr/ecommerce-api/internal/models"
)

func (r *Repo) CreateUser(ctx context.Context, u models.User) error {
	_, err := r.pool.Exec(ctx, `INSERT INTO users (email, password_hash, first_name, last_name, role, is_active)
	VALUES ($1, $2, $3, $4, $5, $6)`,
		u.Email, u.PasswordHash, u.FirstName, u.LastName, u.Role, u.IsActive,
	)

	return err
}

func (r *Repo) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user *models.User
	row := r.pool.QueryRow(ctx, `SELECT id, email, password_hash, first_name, last_name, role, is_active, created_at, updated_at FROM users`)
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.Role,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}
