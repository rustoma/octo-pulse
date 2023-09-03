package storage

import (
	"github.com/rustoma/octo-pulse/internal/models"
)

type Store struct {
	User UserStore
}

type UserStore interface {
	GetUserByEmail(email string) (*models.User, error)
	UpdateRefreshToken(userId int, refreshToken string) (int, error)
	SelectUserByRefreshToken(refreshToken string) (*models.User, error)
	UpdateUserRefreshToken(userId int, refreshToken string) (int, error)
	InsertUser(user *models.User) (int, error)
}
