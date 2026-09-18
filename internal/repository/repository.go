package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/xenptr/ecommerce-api/internal/models"
)

type UserRepository interface {
	CreateUser(context.Context, models.User) error
	GetUserByEmail(context.Context, string) (*models.User, error)
}

type Store interface {
	UserRepository
}

type Repo struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Repo {
	return &Repo{
		pool: pool,
	}
}
