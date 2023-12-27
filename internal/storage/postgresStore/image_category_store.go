package postgresstore

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rustoma/octo-pulse/internal/models"
)

type PostgresImageCategoryStore struct {
	DB        *pgxpool.Pool
	dbTimeout time.Duration
}

func NewImageCategoryStore(DB *pgxpool.Pool) *PostgresImageCategoryStore {
	return &PostgresImageCategoryStore{
		DB:        DB,
		dbTimeout: time.Second * 20,
	}
}

func (a *PostgresImageCategoryStore) InsertCategory(category *models.ImageCategory) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), a.dbTimeout)
	defer cancel()

	stmt, args, err := pgQb().
		Insert("public.image_category").
		Columns("name, created_at, updated_at").
		Values(category.Name, time.Now().UTC(), time.Now().UTC()).
		Suffix("RETURNING \"id\"").
		ToSql()

	if err != nil {
		logger.Err(err).Send()
		return 0, err
	}

	var categoryId int

	err = a.DB.QueryRow(ctx, stmt, args...).Scan(&categoryId)
	return categoryId, err
}
