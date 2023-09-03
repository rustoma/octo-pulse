package postgresstore

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rustoma/octo-pulse/internal/models"
)

type PostgressDomainStore struct {
	DB        *pgxpool.Pool
	dbTimeout time.Duration
}

func NewDomainStore(DB *pgxpool.Pool) *PostgressDomainStore {
	return &PostgressDomainStore{
		DB:        DB,
		dbTimeout: time.Second * 3,
	}
}

func (d *PostgressDomainStore) InsertDomain(domain *models.Domain) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), d.dbTimeout)
	defer cancel()

	stmt, args, err := pgQb().
		Insert("public.domain").
		Columns("name, created_at, updated_at").
		Values(domain.Name, time.Now().UTC(), time.Now().UTC()).
		Suffix("RETURNING \"id\"").
		ToSql()

	if err != nil {
		logger.Err(err).Send()
		return 0, err
	}

	var domainId int

	err = d.DB.QueryRow(ctx, stmt, args...).Scan(&domainId)
	return domainId, err
}
