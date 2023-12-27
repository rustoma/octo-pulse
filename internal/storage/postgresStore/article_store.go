package postgresstore

import (
	"context"
	"github.com/rustoma/octo-pulse/internal/dto"
	"github.com/rustoma/octo-pulse/internal/storage"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rustoma/octo-pulse/internal/models"
)

type PostgressArticleStore struct {
	DB                *pgxpool.Pool
	categoryStore     storage.CategoryStore
	imageStorageStore storage.ImageStorageStore
	authorStore       storage.AuthorStore
	dbTimeout         time.Duration
}

func NewArticleStore(DB *pgxpool.Pool, categoryStore storage.CategoryStore, imageStorageStore storage.ImageStorageStore, authorStore storage.AuthorStore) *PostgressArticleStore {
	return &PostgressArticleStore{
		DB:                DB,
		categoryStore:     categoryStore,
		imageStorageStore: imageStorageStore,
		authorStore:       authorStore,
		dbTimeout:         time.Second * 20,
	}
}

func (s *PostgressArticleStore) InsertArticle(article *models.Article) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.dbTimeout)
	defer cancel()

	stmt, args, err := pgQb().
		Insert("public.article").
		Columns("title, slug, body, thumbnail, publication_date, is_published, author_id, category_id, domain_id, featured, reading_time, is_sponsored,created_at, updated_at").
		Values(article.Title, article.Slug, article.Body, article.Thumbnail, article.PublicationDate, article.IsPublished,
			article.AuthorId, article.CategoryId, article.DomainId, article.Featured, article.ReadingTime, article.IsSponsored, time.Now().UTC(), time.Now().UTC()).
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

func (s *PostgressArticleStore) GetArticles(filters ...*storage.GetArticlesFilters) ([]*dto.Article, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.dbTimeout)
	defer cancel()

	selectStmt := "*"

	if len(filters) > 0 && filters[0].ExcludeBody == "true" {
		selectStmt = "id, title, slug, thumbnail, publication_date, is_published, author_id, category_id, domain_id, featured, reading_time, is_sponsored,created_at, updated_at"
	}

	articlesStmt := pgQb().
		Select(selectStmt).
		OrderBy("created_at DESC").
		From("public.article")

	if len(filters) > 0 && filters[0].Limit != 0 {
		articlesStmt = articlesStmt.Limit(uint64(filters[0].Limit))
	}

	if len(filters) > 0 && filters[0].Offset != 0 {
		articlesStmt = articlesStmt.Offset(uint64(filters[0].Offset))
	}

	if len(filters) > 0 && filters[0].CategoryId != 0 {
		articlesStmt = articlesStmt.Where(
			squirrel.And{
				squirrel.Eq{"category_id": filters[0].CategoryId},
			})
	}

	if len(filters) > 0 && filters[0].DomainId != 0 {
		articlesStmt = articlesStmt.Where(
			squirrel.And{
				squirrel.Eq{"domain_id": filters[0].DomainId},
			})
	}

	if len(filters) > 0 && filters[0].Slug != "" {
		articlesStmt = articlesStmt.Where(
			squirrel.And{
				squirrel.Eq{"slug": filters[0].Slug},
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

	var articles []*dto.Article

	for rows.Next() {
		var articleFromScan *models.Article

		if len(filters) > 0 && filters[0].ExcludeBody == "true" {
			article, err := scanToArticleWithoutBody(rows)
			if err != nil {
				logger.Err(err).Send()
				return nil, err
			}
			articleFromScan = article
		} else {
			article, err := scanToArticle(rows)
			if err != nil {
				logger.Err(err).Send()
				return nil, err
			}
			articleFromScan = article
		}

		dtoArticle := dto.Article{
			ID:              articleFromScan.ID,
			Title:           articleFromScan.Title,
			Slug:            articleFromScan.Slug,
			Body:            articleFromScan.Body,
			PublicationDate: articleFromScan.PublicationDate,
			IsPublished:     articleFromScan.IsPublished,
			DomainId:        articleFromScan.DomainId,
			Featured:        articleFromScan.Featured,
			ReadingTime:     articleFromScan.ReadingTime,
			IsSponsored:     articleFromScan.IsSponsored,
			CreatedAt:       articleFromScan.CreatedAt,
			UpdatedAt:       articleFromScan.UpdatedAt,
		}

		if articleFromScan.Thumbnail != nil {
			thumbnail, err := s.imageStorageStore.GetImage(*articleFromScan.Thumbnail)
			if err != nil {
				logger.Err(err).Send()
				return nil, err
			}
			dtoArticle.Thumbnail = thumbnail
		}

		category, err := s.categoryStore.GetCategory(articleFromScan.CategoryId)
		if err != nil {
			logger.Err(err).Send()
			return nil, err
		}

		if category == nil {
			logger.Error().Msgf("No categories for article %s", articleFromScan.Title)
			return nil, err
		}

		dtoArticle.Category = *category

		author, err := s.authorStore.GetAuthor(articleFromScan.AuthorId)
		if err != nil {
			logger.Err(err).Send()
			return nil, err
		}

		dtoArticle.Author = *author

		if err != nil {
			logger.Err(err).Send()
			return nil, err
		}

		articles = append(articles, &dtoArticle)
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
		&article.Body,
		&article.Thumbnail,
		&article.PublicationDate,
		&article.IsPublished,
		&article.AuthorId,
		&article.CategoryId,
		&article.DomainId,
		&article.Featured,
		&article.ReadingTime,
		&article.IsSponsored,
		&article.CreatedAt,
		&article.UpdatedAt,
	)

	return &article, err
}

func scanToArticleWithoutBody(rows pgx.Rows) (*models.Article, error) {
	var article models.Article
	err := rows.Scan(
		&article.ID,
		&article.Title,
		&article.Slug,
		&article.Thumbnail,
		&article.PublicationDate,
		&article.IsPublished,
		&article.AuthorId,
		&article.CategoryId,
		&article.DomainId,
		&article.Featured,
		&article.ReadingTime,
		&article.IsSponsored,
		&article.CreatedAt,
		&article.UpdatedAt,
	)

	return &article, err
}

func convertArticleToArticleMap(article *models.Article) map[string]interface{} {
	return map[string]interface{}{
		"title":            article.Title,
		"body":             article.Body,
		"slug":             article.Slug,
		"thumbnail":        article.Thumbnail,
		"publication_date": article.PublicationDate,
		"is_published":     article.IsPublished,
		"author_id":        article.AuthorId,
		"category_id":      article.CategoryId,
		"domain_id":        article.DomainId,
		"featured":         article.Featured,
		"reading_time":     article.ReadingTime,
		"is_sponsored":     article.IsSponsored,
		"created_at":       article.CreatedAt,
		"updated_at":       article.UpdatedAt,
	}
}
