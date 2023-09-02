package services

import (
	"context"

	"github.com/rustoma/octo-pulse/internal/models"
	"github.com/rustoma/octo-pulse/internal/storage"
)

type userRoles struct {
	Admin  int
	Editor int
}

type UserService interface {
	GetUserByID(context.Context, int) (*models.User, error)
}

type userService struct {
	store     storage.Store
	userRoles userRoles
}

func NewUserService(store storage.Store) UserService {
	return &userService{store: store, userRoles: userRoles{
		Admin:  1,
		Editor: 2,
	}}
}

func (u *userService) GetUserByID(context.Context, int) (*models.User, error) {
	return &models.User{ID: 1}, nil
}
