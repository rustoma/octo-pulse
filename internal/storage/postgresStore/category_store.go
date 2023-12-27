package postgresstore

import (
	"context"
	"github.com/rustoma/octo-pulse/internal/storage"
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
		dbTimeout: time.Second * 20,
	}
}

func (c *PostgressCategoryStore) InsertCategory(category *models.Category) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.dbTimeout)
	defer cancel()

	stmt, args, err := pgQb().
		Insert("public.category").
		Columns("name, slug, weight, created_at, updated_at").
		Values(category.Name, category.Slug, category.Weight, time.Now().UTC(), time.Now().UTC()).
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

func (s *PostgressCategoryStore) GetCategories(filters ...*storage.GetCategoriesFilters) ([]*models.Category, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.dbTimeout)
	defer cancel()

	categoriesStmt := pgQb().
		Select("*").
		OrderBy("weight").
		OrderBy("name").
		From("public.category")

	if len(filters) > 0 && filters[0].Slug != "" {
		categoriesStmt = categoriesStmt.Where(
			squirrel.And{
				squirrel.Eq{"slug": filters[0].Slug},
			})
	}

	stmt, args, err := categoriesStmt.ToSql()

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
		&category.Slug,
		&category.Weight,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	return &category, err
}
