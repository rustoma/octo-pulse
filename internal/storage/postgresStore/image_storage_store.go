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

type PostgresImageStorageStore struct {
	DB        *pgxpool.Pool
	dbTimeout time.Duration
}

func NewImageStorageStore(DB *pgxpool.Pool) *PostgresImageStorageStore {
	return &PostgresImageStorageStore{
		DB:        DB,
		dbTimeout: time.Second * 3,
	}
}

func (s *PostgresImageStorageStore) InsertImage(image *models.Image) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.dbTimeout)
	defer cancel()

	stmt, args, err := pgQb().
		Insert("public.image_storage").
		Columns("name, path, size, type, width, height, alt, category_id, created_at, updated_at").
		Values(image.Name, image.Path, image.Size, image.Type, image.Width, image.Height, image.Alt, image.CategoryId, time.Now().UTC(), time.Now().UTC()).
		Suffix("RETURNING \"id\"").
		ToSql()

	if err != nil {
		logger.Err(err).Send()
		return 0, err
	}

	var imageId int

	err = s.DB.QueryRow(ctx, stmt, args...).Scan(&imageId)
	return imageId, err
}

func (s *PostgresImageStorageStore) GetImage(id int) (*models.Image, error) {

	ctx, cancel := context.WithTimeout(context.Background(), s.dbTimeout)
	defer cancel()

	stmt, args, err := pgQb().
		Select("*").
		From("public.image_storage").
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

	var image *models.Image

	for rows.Next() {
		imageFromScan, err := scanToImage(rows)

		if err != nil {
			logger.Err(err).Send()
			return nil, err
		}

		image = imageFromScan
	}

	return image, err
}

func (s *PostgresImageStorageStore) GetImages(filters ...*storage.GetImagesFilters) ([]*models.Image, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.dbTimeout)
	defer cancel()

	imagesStmt := pgQb().
		Select("*").
		From("public.image_storage")

	if len(filters) > 0 && filters[0].CategoryId != 0 {
		imagesStmt = imagesStmt.Where(
			squirrel.And{
				squirrel.Eq{"category_id": filters[0].CategoryId},
			})
	}

	stmt, args, err := imagesStmt.ToSql()

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

	var images []*models.Image

	for rows.Next() {
		imageFromScan, err := scanToImage(rows)

		if err != nil {
			logger.Err(err).Send()
			return nil, err
		}

		images = append(images, imageFromScan)
	}

	return images, err
}

func scanToImage(rows pgx.Rows) (*models.Image, error) {
	var image models.Image
	err := rows.Scan(
		&image.ID,
		&image.Name,
		&image.Path,
		&image.Size,
		&image.Type,
		&image.Width,
		&image.Height,
		&image.Alt,
		&image.CategoryId,
		&image.CreatedAt,
		&image.UpdatedAt,
	)

	return &image, err
}
