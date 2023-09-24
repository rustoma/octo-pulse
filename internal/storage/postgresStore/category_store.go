package postgresstore

import (
	"context"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
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

func (s *PostgressCategoryStore) GetCategories() ([]*models.Category, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.dbTimeout)
	defer cancel()

	stmt, args, err := pgQb().
		Select("*").
		From("public.category").
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

	var categories []*models.Category

	for rows.Next() {
		categoryFromScan, err := scanToCategory(rows)

		if err != nil {
			logger.Err(err).Send()
			return nil, err
		}

		categories = append(categories, categoryFromScan)
	}

	return categories, err
}

func (s *PostgressCategoryStore) GetCategory(id int) (*models.Category, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.dbTimeout)
	defer cancel()

	stmt, args, err := pgQb().
		Select("*").
		From("public.category").
		Where(squirrel.Eq{"id": id}).
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

	var category *models.Category

	for rows.Next() {
		categoryFromScan, err := scanToCategory(rows)

		if err != nil {
			logger.Err(err).Send()
			return nil, err
		}

		category = categoryFromScan
	}

	return category, err
}

func scanToCategory(rows pgx.Rows) (*models.Category, error) {
	var category models.Category
	err := rows.Scan(
		&category.ID,
		&category.Name,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	return &category, err
}
