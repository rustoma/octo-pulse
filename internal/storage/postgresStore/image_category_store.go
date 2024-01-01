package postgresstore

import (
	"context"
	"github.com/jackc/pgx/v5"
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

func (s *PostgresImageCategoryStore) InsertCategory(category *models.ImageCategory) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.dbTimeout)
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

	err = s.DB.QueryRow(ctx, stmt, args...).Scan(&categoryId)
	return categoryId, err
}

func (s *PostgresImageCategoryStore) GetCategories() ([]*models.ImageCategory, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.dbTimeout)
	defer cancel()

	stmt, args, err := pgQb().
		Select("*").
		OrderBy("name ASC").
		From("public.image_category").
		ToSql()

	if err != nil {
		logger.Err(err).Send()
		return nil, err
	}

	rows, err := s.DB.Query(ctx, stmt, args...)
	defer rows.Close()

	if err != nil {
		logger.Err(err).Send()
		return nil, err
	}

	categories := make([]*models.ImageCategory, 0)

	for rows.Next() {
		categoryFromScan, err := scanToImageCategory(rows)

		if err != nil {
			logger.Err(err).Send()
			return nil, err
		}

		categories = append(categories, categoryFromScan)
	}

	return categories, err
}

func scanToImageCategory(rows pgx.Rows) (*models.ImageCategory, error) {
	var imageCategory models.ImageCategory

	err := rows.Scan(
		&imageCategory.ID,
		&imageCategory.Name,
		&imageCategory.CreatedAt,
		&imageCategory.UpdatedAt,
	)

	return &imageCategory, err
}
