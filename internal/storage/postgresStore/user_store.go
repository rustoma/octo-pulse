package postgresstore

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rustoma/octo-pulse/internal/models"
)

type PostgressUserStore struct {
	DB *pgxpool.Pool
}

func NewUserStore(DB *pgxpool.Pool) *PostgressUserStore {
	return &PostgressUserStore{
		DB: DB,
	}
}

func (u *PostgressUserStore) GetUserByID(context.Context, int) (*models.User, error) {
	return &models.User{ID: 1}, nil
}
