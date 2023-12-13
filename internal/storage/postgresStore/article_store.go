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

type PostgressArticleStore struct {
	DB        *pgxpool.Pool
	dbTimeout time.Duration
}

func NewArticleStore(DB *pgxpool.Pool) *PostgressArticleStore {
	return &PostgressArticleStore{
		DB:        DB,
		dbTimeout: time.Second * 3,
	}
}

func (s *PostgressArticleStore) InsertArticle(article *models.Article) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.dbTimeout)
	defer cancel()

	stmt, args, err := pgQb().
		Insert("public.article").
		Columns("title, slug, description, thumbnail, publication_date, is_published, author_id, category_id, domain_id, featured, created_at, updated_at").
		Values(article.Title, article.Slug, article.Description, article.Thumbnail, article.PublicationDate, article.IsPublished,
			article.AuthorId, article.CategoryId, article.DomainId, article.Featured, time.Now().UTC(), time.Now().UTC()).
		Suffix("RETURNING \"id\"").
		ToSql()

	if err != nil {
		logger.Err(err).Send()
		return 0, err
	}

	var articleId int

	err = s.DB.QueryRow(ctx, stmt, args...).Scan(&articleId)
	return articleId, err
}

func (s *PostgressArticleStore) DeleteArticle(id int) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.dbTimeout)
	defer cancel()

	stmt, args, err := pgQb().
		Delete("public.article").
		Where(squirrel.Eq{"id": id}).
		Suffix("RETURNING \"id\"").
		ToSql()

	if err != nil {
		return 0, err
	}

	var articleId int

	err = s.DB.QueryRow(ctx, stmt, args...).Scan(&articleId)
	return articleId, err
}

func (s *PostgressArticleStore) GetArticles(filters ...*storage.GetArticlesFilters) ([]*models.Article, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.dbTimeout)
	defer cancel()

	articlesStmt := pgQb().
		Select("*").
		OrderBy("created_at DESC").
		From("public.article")

	if len(filters) > 0 && filters[0].Limit != 0 {
		articlesStmt = articlesStmt.Limit(uint64(filters[0].Limit))
	}

	if len(filters) > 0 && filters[0].CategoryId != 0 {
		articlesStmt = articlesStmt.Where(
			squirrel.And{
				squirrel.Eq{"category_id": filters[0].CategoryId},
			})
	}

	if len(filters) > 0 && (filters[0].Featured == "true" || filters[0].Featured == "false") {
		featured := false

		if filters[0].Featured == "true" {
			featured = true
		}

		articlesStmt = articlesStmt.Where(
			squirrel.And{
				squirrel.Eq{"featured": featured},
			})
	}

	stmt, args, err := articlesStmt.ToSql()

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

	var articles []*models.Article

	for rows.Next() {
		articleFromScan, err := scanToArticle(rows)

		if err != nil {
			logger.Err(err).Send()
			return nil, err
		}

		articles = append(articles, articleFromScan)
	}

	return articles, err

}

func (s *PostgressArticleStore) GetArticle(id int) (*models.Article, error) {

	ctx, cancel := context.WithTimeout(context.Background(), s.dbTimeout)
	defer cancel()

	stmt, args, err := pgQb().
		Select("*").
		From("public.article").
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

	var article *models.Article

	for rows.Next() {
		articleFromScan, err := scanToArticle(rows)

		if err != nil {
			logger.Err(err).Send()
			return nil, err
		}

		article = articleFromScan
	}

	return article, err
}

func (s *PostgressArticleStore) UpdateArticle(id int, article *models.Article) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.dbTimeout)
	defer cancel()

	articleMap := convertArticleToArticleMap(article)
	articleMap["updated_at"] = time.Now().UTC()

	stmt, args, err := pgQb().
		Update("public.article").
		SetMap(articleMap).
		Where(squirrel.Eq{"id": id}).
		Suffix("RETURNING \"id\"").ToSql()

	if err != nil {
		logger.Err(err).Send()
		return 0, err
	}

	var updatedArticleId int

	err = s.DB.QueryRow(ctx, stmt, args...).Scan(&updatedArticleId)
	return updatedArticleId, err
}

func scanToArticle(rows pgx.Rows) (*models.Article, error) {
	var article models.Article
	err := rows.Scan(
		&article.ID,
		&article.Title,
		&article.Slug,
		&article.Description,
		&article.Thumbnail,
		&article.PublicationDate,
		&article.IsPublished,
		&article.AuthorId,
		&article.CategoryId,
		&article.DomainId,
		&article.Featured,
		&article.CreatedAt,
		&article.UpdatedAt,
	)

	return &article, err
}

func convertArticleToArticleMap(article *models.Article) map[string]interface{} {
	return map[string]interface{}{
		"title":            article.Title,
		"description":      article.Description,
		"thumbnail":        article.Thumbnail,
		"publication_date": article.PublicationDate,
		"is_published":     article.IsPublished,
		"author_id":        article.AuthorId,
		"category_id":      article.CategoryId,
		"domain_id":        article.DomainId,
		"created_at":       article.CreatedAt,
		"updated_at":       article.UpdatedAt,
	}
}
