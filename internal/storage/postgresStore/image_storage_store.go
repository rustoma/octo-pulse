package postgresstore

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rustoma/octo-pulse/internal/models"
)

type PostgresImageStorageStore struct {
	DB        *pgxpool.Pool
	dbTimeout time.Duration
}

func NewImageStorageStore(DB *pgxpool.Pool) *PostgresImageStorageStore {
	return &PostgresImageStorageStore{
		DB:        DB,
		dbTimeout: time.Second * 3,
	}
}

func (a *PostgresImageStorageStore) InsertImage(image *models.Image) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), a.dbTimeout)
	defer cancel()

	stmt, args, err := pgQb().
		Insert("public.image_storage").
		Columns("name, path, size, type, width, height, alt, category_id, created_at, updated_at").
		Values(image.Name, image.Path, image.Size, image.Type, image.Width, image.Height, image.Alt, image.CategoryId, time.Now().UTC(), time.Now().UTC()).
		Suffix("RETURNING \"id\"").
		ToSql()

	if err != nil {
		logger.Err(err).Send()
		return 0, err
	}

	var imageId int

	err = a.DB.QueryRow(ctx, stmt, args...).Scan(&imageId)
	return imageId, err
}
