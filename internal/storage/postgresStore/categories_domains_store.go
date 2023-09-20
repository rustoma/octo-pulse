package postgresstore

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgressCategoriesDomainsStore struct {
	DB        *pgxpool.Pool
	dbTimeout time.Duration
}

func NewCategoriesDomainsStore(DB *pgxpool.Pool) *PostgressCategoriesDomainsStore {
	return &PostgressCategoriesDomainsStore{
		DB:        DB,
		dbTimeout: time.Second * 3,
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
