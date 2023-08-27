package postgresstore

import "github.com/rustoma/octo-pulse/internal/storage"

func NewPostgresStorage() *storage.Store {
	return &storage.Store{
		User: newUserStore(),
	}
}
