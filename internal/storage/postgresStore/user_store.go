package postgresstore

import (
	"context"

	"github.com/rustoma/octo-pulse/internal/model"
)

type PostgressUserStore struct{}

func newUserStore() *PostgressUserStore {
	return &PostgressUserStore{}
}

func (u *PostgressUserStore) GetUserByID(context.Context, int) (*model.User, error) {
	return &model.User{ID: 1}, nil
}
