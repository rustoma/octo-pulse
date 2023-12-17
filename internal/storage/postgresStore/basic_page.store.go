package postgresstore

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/rustoma/octo-pulse/internal/storage"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rustoma/octo-pulse/internal/models"
)

type PostgresBasicPageStore struct {
	DB        *pgxpool.Pool
	dbTimeout time.Duration
}

func NewBasicPageStore(DB *pgxpool.Pool) *PostgresBasicPageStore {
	return &PostgresBasicPageStore{
		DB:        DB,
		dbTimeout: time.Second * 3,
	}
}

func (s *PostgresBasicPageStore) InsertBasicPage(page *models.BasicPage) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.dbTimeout)
	defer cancel()

	stmt, args, err := pgQb().
		Insert("public.basic_page").
		Columns("title, slug, body, domain, created_at, updated_at").
		Values(page.Title, page.Slug, page.Body, page.Domain, time.Now().UTC(), time.Now().UTC()).
		Suffix("RETURNING \"id\"").
		ToSql()

	if err != nil {
		logger.Err(err).Send()
		return 0, err
	}

	var pageId int

	err = s.DB.QueryRow(ctx, stmt, args...).Scan(&pageId)
	return pageId, err
}

func (s *PostgresBasicPageStore) GetBasicPages(filters ...*storage.GetBasicPagesFilters) ([]*models.BasicPage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.dbTimeout)
	defer cancel()

	basicPagesStmt := pgQb().
		Select("*").
		OrderBy("created_at DESC").
		From("public.basic_page")

	if len(filters) > 0 && filters[0].DomainId != 0 {
		basicPagesStmt = basicPagesStmt.Where(
			squirrel.And{
				squirrel.Eq{"domain": filters[0].DomainId},
			})
	}

	stmt, args, err := basicPagesStmt.ToSql()

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

	var basicPages []*models.BasicPage

	for rows.Next() {
		basicPagesFromScan, err := scanToBasicPage(rows)

		if err != nil {
			logger.Err(err).Send()
			return nil, err
		}

		basicPages = append(basicPages, basicPagesFromScan)
	}

	return basicPages, err
}

func (s *PostgresBasicPageStore) GetBasicPage(id int) (*models.BasicPage, error) {

	ctx, cancel := context.WithTimeout(context.Background(), s.dbTimeout)
	defer cancel()

	stmt, args, err := pgQb().
		Select("*").
		From("public.basic_page").
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

	var page *models.BasicPage

	for rows.Next() {
		pageFromScan, err := scanToBasicPage(rows)

		if err != nil {
			logger.Err(err).Send()
			return nil, err
		}

		page = pageFromScan
	}

	return page, err
}

func (s *PostgresBasicPageStore) GetBasicPageBySlug(slug string, filters ...*storage.GetBasicPageBySlugFilters) (*models.BasicPage, error) {

	ctx, cancel := context.WithTimeout(context.Background(), s.dbTimeout)
	defer cancel()

	getBasicPageBySlugstmt := pgQb().
		Select("*").
		From("public.basic_page").
		Where(squirrel.Eq{"slug": slug})

	if len(filters) > 0 && filters[0].DomainId != 0 {
		getBasicPageBySlugstmt = getBasicPageBySlugstmt.Where(
			squirrel.And{
				squirrel.Eq{"domain": filters[0].DomainId},
			})
	}

	stmt, args, err := getBasicPageBySlugstmt.ToSql()

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

	var page *models.BasicPage

	for rows.Next() {
		pageFromScan, err := scanToBasicPage(rows)

		if err != nil {
			logger.Err(err).Send()
			return nil, err
		}

		page = pageFromScan
	}

	return page, err
}

func scanToBasicPage(rows pgx.Rows) (*models.BasicPage, error) {
	var page models.BasicPage
	err := rows.Scan(
		&page.ID,
		&page.Title,
		&page.Slug,
		&page.Body,
		&page.Domain,
		&page.CreatedAt,
		&page.UpdatedAt,
	)

	return &page, err
}
