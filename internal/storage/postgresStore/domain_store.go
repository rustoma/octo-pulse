package postgresstore

import (
	"context"
	"github.com/rustoma/octo-pulse/internal/dto"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rustoma/octo-pulse/internal/models"
)

type PostgresDomainStore struct {
	DB        *pgxpool.Pool
	dbTimeout time.Duration
}

func NewDomainStore(DB *pgxpool.Pool) *PostgresDomainStore {
	return &PostgresDomainStore{
		DB:        DB,
		dbTimeout: time.Second * 20,
	}
}

func (s *PostgresDomainStore) InsertDomain(domain *models.Domain) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.dbTimeout)
	defer cancel()

	stmt, args, err := pgQb().
		Insert("public.domain").
		Columns("name, email, created_at, updated_at").
		Values(domain.Name, domain.Email, time.Now().UTC(), time.Now().UTC()).
		Suffix("RETURNING \"id\"").
		ToSql()

	if err != nil {
		logger.Err(err).Send()
		return 0, err
	}

	var domainId int

	err = s.DB.QueryRow(ctx, stmt, args...).Scan(&domainId)
	return domainId, err
}

func (s *PostgresDomainStore) GetDomains() ([]*models.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.dbTimeout)
	defer cancel()

	stmt, args, err := pgQb().
		Select("*").
		From("public.domain").
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

	var domains []*models.Domain

	for rows.Next() {
		domainFromScan, err := scanToDomain(rows)

		if err != nil {
			logger.Err(err).Send()
			return nil, err
		}

		domains = append(domains, domainFromScan)
	}

	return domains, err
}

func (s *PostgresDomainStore) GetDomain(id int) (*models.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.dbTimeout)
	defer cancel()

	stmt, args, err := pgQb().
		Select("*").
		From("public.domain").
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

	var domain *models.Domain

	for rows.Next() {
		domainFromScan, err := scanToDomain(rows)

		if err != nil {
			logger.Err(err).Send()
			return nil, err
		}

		domain = domainFromScan
	}

	return domain, err
}

func (s *PostgresDomainStore) GetDomainPublicData(id int) (*dto.DomainPublicData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.dbTimeout)
	defer cancel()

	stmt, args, err := pgQb().
		Select("email").
		From("public.domain").
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

	var domainPublicData *dto.DomainPublicData

	for rows.Next() {
		var domainPublicDataFromScan dto.DomainPublicData

		err := rows.Scan(
			&domainPublicDataFromScan.Email,
		)

		if err != nil {
			logger.Err(err).Send()
			return nil, err
		}

		domainPublicData = &domainPublicDataFromScan
	}

	return domainPublicData, err
}

func (s *PostgresDomainStore) UpdateDomain(id int, domain *models.Domain) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.dbTimeout)
	defer cancel()

	domainMap := convertDomainToDomainMap(domain)
	domainMap["updated_at"] = time.Now().UTC()

	stmt, args, err := pgQb().
		Update("public.domain").
		SetMap(domainMap).
		Where(squirrel.Eq{"id": id}).
		Suffix("RETURNING \"id\"").ToSql()

	if err != nil {
		logger.Err(err).Send()
		return 0, err
	}

	var updatedDomainId int

	err = s.DB.QueryRow(ctx, stmt, args...).Scan(&updatedDomainId)
	return updatedDomainId, err
}

func scanToDomain(rows pgx.Rows) (*models.Domain, error) {
	var domain models.Domain
	err := rows.Scan(
		&domain.ID,
		&domain.Name,
		&domain.Email,
		&domain.CreatedAt,
		&domain.UpdatedAt,
	)

	return &domain, err
}

func convertDomainToDomainMap(domain *models.Domain) map[string]interface{} {
	return map[string]interface{}{
		"name":       domain.Name,
		"email":      domain.Email,
		"created_at": domain.CreatedAt,
		"updated_at": domain.UpdatedAt,
	}
}
