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

type PostgresCategoryStore struct {
	DB        *pgxpool.Pool
	dbTimeout time.Duration
}

func NewCategoryStore(DB *pgxpool.Pool) *PostgresCategoryStore {
	return &PostgresCategoryStore{
		DB:        DB,
		dbTimeout: time.Second * 20,
	}
}

func (c *PostgresCategoryStore) InsertCategory(category *models.Category) (int, error) {
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

func (s *PostgresCategoryStore) GetCategories(filters ...*storage.GetCategoriesFilters) ([]*models.Category, error) {
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

func (s *PostgresCategoryStore) GetCategory(id int) (*models.Category, error) {
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

func (s *PostgresCategoryStore) UpdateCategory(id int, category *models.Category) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.dbTimeout)
	defer cancel()

	categoryMap := convertCategoryToCategoryMap(category)
	categoryMap["updated_at"] = time.Now().UTC()

	stmt, args, err := pgQb().
		Update("public.category").
		SetMap(categoryMap).
		Where(squirrel.Eq{"id": id}).
		Suffix("RETURNING \"id\"").ToSql()

	if err != nil {
		logger.Err(err).Send()
		return 0, err
	}

	var updatedCategoryId int

	err = s.DB.QueryRow(ctx, stmt, args...).Scan(&updatedCategoryId)
	return updatedCategoryId, err
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

func convertCategoryToCategoryMap(category *models.Category) map[string]interface{} {
	return map[string]interface{}{
		"name":       category.Name,
		"slug":       category.Slug,
		"weight":     category.Weight,
		"created_at": category.CreatedAt,
		"updated_at": category.UpdatedAt,
	}
}
