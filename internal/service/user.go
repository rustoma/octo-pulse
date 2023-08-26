package service

import (
	"context"

	"github.com/rustoma/octo-pulse/internal/model"
	"github.com/rustoma/octo-pulse/internal/storage"
)

type UserService interface {
	GetUserByID(context.Context, int) (*model.User, error)
}

type userService struct {
	store storage.Store
}

func NewAnswerService(store storage.Store) UserService {
	return &userService{store: store}
}

func (u *userService) GetUserByID(context.Context, int) (*model.User, error) {
	return &model.User{ID: 1}, nil
}
