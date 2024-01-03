package postgresstore

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rustoma/octo-pulse/internal/models"
)

type PostgresAuthorStore struct {
	DB        *pgxpool.Pool
	dbTimeout time.Duration
}

func NewAuthorStore(DB *pgxpool.Pool) *PostgresAuthorStore {
	return &PostgresAuthorStore{
		DB:        DB,
		dbTimeout: time.Second * 20,
	}
}

func (s *PostgresAuthorStore) InsertAuthor(author *models.Author) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.dbTimeout)
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

	err = s.DB.QueryRow(ctx, stmt, args...).Scan(&authorId)
	return authorId, err
}

func (s *PostgresAuthorStore) GetAuthors() ([]*models.Author, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.dbTimeout)
	defer cancel()

	stmt, args, err := pgQb().
		Select("*").
		OrderBy("first_name ASC").
		From("public.author").
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

	var authors []*models.Author

	for rows.Next() {
		authorsFromScan, err := scanToAuthor(rows)

		if err != nil {
			logger.Err(err).Send()
			return nil, err
		}

		authors = append(authors, authorsFromScan)
	}

	return authors, err
}

func (s *PostgresAuthorStore) GetAuthor(id int) (*models.Author, error) {

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

func (s *PostgresAuthorStore) UpdateAuthor(id int, author *models.Author) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.dbTimeout)
	defer cancel()

	authorMap := convertAuthorToAuthorMap(author)
	authorMap["updated_at"] = time.Now().UTC()

	stmt, args, err := pgQb().
		Update("public.author").
		SetMap(authorMap).
		Where(squirrel.Eq{"id": id}).
		Suffix("RETURNING \"id\"").ToSql()

	if err != nil {
		logger.Err(err).Send()
		return 0, err
	}

	var updatedAuthorId int

	err = s.DB.QueryRow(ctx, stmt, args...).Scan(&updatedAuthorId)
	return updatedAuthorId, err
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

func convertAuthorToAuthorMap(author *models.Author) map[string]interface{} {
	return map[string]interface{}{
		"first_name":  author.FirstName,
		"last_name":   author.LastName,
		"description": author.Description,
		"image_url":   author.ImageUrl,
		"created_at":  author.CreatedAt,
		"updated_at":  author.UpdatedAt,
	}
}
