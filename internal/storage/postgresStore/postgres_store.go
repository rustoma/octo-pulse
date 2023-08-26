package postgresstore

import "github.com/rustoma/octo-pulse/internal/storage"

type store struct{}

func newPostgresStorage() *storage.Store {
	return &store{}
}
