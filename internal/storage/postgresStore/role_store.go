package postgresstore

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rustoma/octo-pulse/internal/models"
)

type PostgressRoleStore struct {
	DB        *pgxpool.Pool
	dbTimeout time.Duration
}

func NewRoleStore(DB *pgxpool.Pool) *PostgressRoleStore {
	return &PostgressRoleStore{
		DB:        DB,
		dbTimeout: time.Second * 20,
	}
}

func (r *PostgressRoleStore) InsertRole(role *models.Role) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.dbTimeout)
	defer cancel()

	stmt, args, err := pgQb().
		Insert("public.role").
		Columns("name, created_at, updated_at").
		Values(role.Name, time.Now().UTC(), time.Now().UTC()).
		Suffix("RETURNING \"id\"").
		ToSql()

	if err != nil {
		logger.Err(err).Send()
		return 0, err
	}

	var roleId int

	err = r.DB.QueryRow(ctx, stmt, args...).Scan(&roleId)
	return roleId, err
}
