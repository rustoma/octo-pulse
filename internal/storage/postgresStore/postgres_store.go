package postgresstore

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rustoma/octo-pulse/internal/storage"
)

func NewPostgresStorage(DB *pgxpool.Pool) *storage.Store {
	return &storage.Store{
		User: NewUserStore(DB),
	}
}
