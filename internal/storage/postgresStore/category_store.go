package postgresstore

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rustoma/octo-pulse/internal/models"
)

type PostgressCategoryStore struct {
	DB        *pgxpool.Pool
	dbTimeout time.Duration
}

func NewCategoryStore(DB *pgxpool.Pool) *PostgressCategoryStore {
	return &PostgressCategoryStore{
		DB:        DB,
		dbTimeout: time.Second * 3,
	}
}

func (c *PostgressCategoryStore) InsertCategory(category *models.Category) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.dbTimeout)
	defer cancel()

	stmt, args, err := pgQb().
		Insert("public.category").
		Columns("name, created_at, updated_at").
		Values(category.Name, time.Now().UTC(), time.Now().UTC()).
		Suffix("RETURNING \"id\"").
		ToSql()

	if err != nil {
		logger.Err(err).Send()
		return 0, err
	}

	var categoryId int

	err = c.DB.QueryRow(ctx, stmt, args...).Scan(&categoryId)
	return categoryId, err
}
