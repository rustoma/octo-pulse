package postgresstore

import (
	"context"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgressCategoriesDomainsStore struct {
	DB        *pgxpool.Pool
	dbTimeout time.Duration
}

func NewCategoriesDomainsStore(DB *pgxpool.Pool) *PostgressCategoriesDomainsStore {
	return &PostgressCategoriesDomainsStore{
		DB:        DB,
		dbTimeout: time.Second * 20,
	}
}

func (s *PostgressCategoriesDomainsStore) AsignCategoryToDomain(categoryId int, domainId int) error {
	ctx, cancel := context.WithTimeout(context.Background(), s.dbTimeout)
	defer cancel()

	stmt, args, err := pgQb().
		Insert("public.categories_domains").
		Columns("domain_id, category_id, created_at").
		Values(domainId, categoryId, time.Now().UTC()).
		Suffix("RETURNING \"category_id\",\"domain_id\"").
		ToSql()

	if err != nil {
		logger.Err(err).Send()
		return err
	}

	var assingedDomainId int
	var assignedCategoryId int

	err = s.DB.QueryRow(ctx, stmt, args...).Scan(&assingedDomainId, &assignedCategoryId)
	return err
}

func (s *PostgressCategoriesDomainsStore) GetDomainCategories(domainId int) ([]int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.dbTimeout)
	defer cancel()

	stmt, args, err := pgQb().
		Select("category_id").
		From("public.categories_domains").
		Where(squirrel.Eq{"domain_id": domainId}).
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

	var categoriesId []int

	for rows.Next() {
		categoryIdFromScan, err := scanToCategoryId(rows)

		if err != nil {
			logger.Err(err).Send()
			return nil, err
		}

		categoriesId = append(categoriesId, categoryIdFromScan)
	}

	return categoriesId, err
}

func scanToCategoryId(rows pgx.Rows) (int, error) {
	var categoryId int
	err := rows.Scan(
		&categoryId,
	)

	return categoryId, err
}
