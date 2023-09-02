package postgresstore

import (
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

func (u *PostgressUserStore) GetUserByEmail(email string) (*models.User, error) {
	return &models.User{ID: 1}, nil
}

func (u *PostgressUserStore) UpdateRefreshToken(userId int, refreshToken string) (int, error) {
	return 1, nil
}

func (u *PostgressUserStore) SelectUserByRefreshToken(refreshToken string) (*models.User, error) {
	return &models.User{ID: 1}, nil
}

func (u *PostgressUserStore) UpdateUserRefreshToken(userId int, refreshToken string) (int, error) {
	return 1, nil
}
