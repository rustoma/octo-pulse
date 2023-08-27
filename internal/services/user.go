package services

import (
	"context"

	"github.com/rustoma/octo-pulse/internal/models"
	"github.com/rustoma/octo-pulse/internal/storage"
)

type UserService interface {
	GetUserByID(context.Context, int) (*models.User, error)
}

type userService struct {
	store storage.Store
}

func NewAnswerService(store storage.Store) UserService {
	return &userService{store: store}
}

func (u *userService) GetUserByID(context.Context, int) (*models.User, error) {
	return &models.User{ID: 1}, nil
}
