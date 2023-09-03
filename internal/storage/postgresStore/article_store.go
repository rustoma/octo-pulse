package postgresstore

import (
	"context"
	"time"

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

func (a *PostgressArticleStore) InsertArticle(article *models.Article) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), a.dbTimeout)
	defer cancel()

	stmt, args, err := pgQb().
		Insert("public.article").
		Columns("title, description, image_url, publication_date, is_published, author_id, category_id, domain_id, created_at, updated_at").
		Values(article.Title, article.Description, article.ImageUrl, article.PublicationDate, article.IsPublished,
			article.AuthorId, article.CategoryId, article.DomainId, time.Now().UTC(), time.Now().UTC()).
		Suffix("RETURNING \"id\"").
		ToSql()

	if err != nil {
		logger.Err(err).Send()
		return 0, err
	}

	var articleId int

	err = a.DB.QueryRow(ctx, stmt, args...).Scan(&articleId)
	return articleId, err
}
