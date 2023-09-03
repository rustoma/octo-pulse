package postgresstore

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rustoma/octo-pulse/internal/models"
)

type PostgressAuthorStore struct {
	DB        *pgxpool.Pool
	dbTimeout time.Duration
}

func NewAuthorStore(DB *pgxpool.Pool) *PostgressAuthorStore {
	return &PostgressAuthorStore{
		DB:        DB,
		dbTimeout: time.Second * 3,
	}
}

func (a *PostgressAuthorStore) InsertAuthor(author *models.Author) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), a.dbTimeout)
	defer cancel()

	stmt, args, err := pgQb().
		Insert("public.author").
		Columns("first_name, last_name, description, image_url, created_at, updated_at").
		Values(author.FirstName, author.LastName, author.Description, author.ImageUrl, time.Now().UTC(), time.Now().UTC()).
		Suffix("RETURNING \"id\"").
		ToSql()

	if err != nil {
		logger.Err(err).Send()
		return 0, err
	}

	var authorId int

	err = a.DB.QueryRow(ctx, stmt, args...).Scan(&authorId)
	return authorId, err
}
