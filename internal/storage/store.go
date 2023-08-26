package storage

import (
	"context"

	"github.com/rustoma/octo-pulse/internal/model"
)

type Store struct {
	User UserStore
}

type UserStore interface {
	GetUserByID(context.Context, int) (*model.User, error)
}
