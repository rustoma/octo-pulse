package postgresstore

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rustoma/octo-pulse/internal/app"
	"github.com/rustoma/octo-pulse/internal/models"
)

type PostgressUserStore struct {
	ctx *app.Ctx
	DB  *pgxpool.Pool
}

func NewUserStore(ctx *app.Ctx, DB *pgxpool.Pool) *PostgressUserStore {
	return &PostgressUserStore{
		ctx: ctx,
		DB:  DB,
	}
}

func (u *PostgressUserStore) GetUserByID(context.Context, int) (*models.User, error) {
	return &models.User{ID: 1}, nil
}
