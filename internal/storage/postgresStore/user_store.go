package postgresstore

import (
	"context"

	"github.com/rustoma/octo-pulse/internal/models"
)

type PostgressUserStore struct{}

func newUserStore() *PostgressUserStore {
	return &PostgressUserStore{}
}

func (u *PostgressUserStore) GetUserByID(context.Context, string) (*models.User, error) {
	return &models.User{ID: 1}, nil
}
