package postgresstore

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
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

func (s *PostgressAuthorStore) GetAuthor(id int) (*models.Author, error) {

	ctx, cancel := context.WithTimeout(context.Background(), s.dbTimeout)
	defer cancel()

	stmt, args, err := pgQb().
		Select("*").
		From("public.author").
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

	var author *models.Author

	for rows.Next() {
		authorFromScan, err := scanToAuthor(rows)

		if err != nil {
			logger.Err(err).Send()
			return nil, err
		}

		author = authorFromScan
	}

	return author, err
}

func scanToAuthor(rows pgx.Rows) (*models.Author, error) {
	var author models.Author
	err := rows.Scan(
		&author.ID,
		&author.FirstName,
		&author.LastName,
		&author.Description,
		&author.ImageUrl,
		&author.CreatedAt,
		&author.UpdatedAt,
	)

	return &author, err
}
