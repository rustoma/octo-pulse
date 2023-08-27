package storage

import (
	"context"

	"github.com/rustoma/octo-pulse/internal/models"
)

type Store struct {
	User UserStore
}

type UserStore interface {
	GetUserByID(context.Context, int) (*models.User, error)
}
