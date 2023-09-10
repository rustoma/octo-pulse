package postgresstore

import (
	"context"
	"errors"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rustoma/octo-pulse/internal/models"
)

type PostgressUserStore struct {
	DB        *pgxpool.Pool
	dbTimeout time.Duration
}

func NewUserStore(DB *pgxpool.Pool) *PostgressUserStore {
	return &PostgressUserStore{
		DB:        DB,
		dbTimeout: time.Second * 3,
	}
}

func (u *PostgressUserStore) InsertUser(user *models.User) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), u.dbTimeout)
	defer cancel()

	stmt, args, err := pgQb().
		Insert("public.user").
		Columns("email, password_hash, role_id, created_at, updated_at").
		Values(user.Email, user.PasswordHash, user.RoleId, time.Now().UTC(), time.Now().UTC()).
		Suffix("RETURNING \"id\"").
		ToSql()

	if err != nil {
		logger.Err(err).Send()
		return 0, err
	}

	var userId int

	err = u.DB.QueryRow(ctx, stmt, args...).Scan(&userId)
	return userId, err

}

func (u *PostgressUserStore) GetUserByEmail(email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), u.dbTimeout)
	defer cancel()

	stmt, args, err := pgQb().
		Select("id, email, COALESCE(refresh_token, '') AS refresh_token, password_hash, role_id, created_at, updated_at").
		From("public.user").
		Where(squirrel.Eq{"email": email}).
		ToSql()

	if err != nil {
		logger.Err(err).Send()
		return nil, err
	}

	rows, err := u.DB.Query(ctx, stmt, args...)
	defer rows.Close()

	if err != nil {
		logger.Err(err).Send()
		return nil, err
	}

	var user *models.User

	for rows.Next() {
		userFromScan, err := scanToUser(rows)

		if err != nil {
			logger.Err(err).Send()
			return nil, err
		}

		user = userFromScan
	}

	if user == nil {
		return nil, errors.New("no user found")
	}

	return user, err
}

func (u *PostgressUserStore) UpdateRefreshToken(userId int, refreshToken string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), u.dbTimeout)
	defer cancel()

	stmt, args, err := pgQb().
		Update("public.user").
		SetMap(map[string]interface{}{
			"refresh_token": refreshToken,
		}).
		Where(squirrel.Eq{"id": userId}).
		Suffix("RETURNING \"id\"").ToSql()

	if err != nil {
		logger.Err(err).Send()
		return 0, err
	}

	var updatedUserId int

	err = u.DB.QueryRow(ctx, stmt, args...).Scan(&updatedUserId)
	return userId, err
}

func (u *PostgressUserStore) SelectUserByRefreshToken(refreshToken string) (*models.User, error) {
	return &models.User{ID: 1}, nil
}

func (u *PostgressUserStore) UpdateUserRefreshToken(userId int, refreshToken string) (int, error) {
	return 1, nil
}

func scanToUser(rows pgx.Rows) (*models.User, error) {
	var user models.User
	err := rows.Scan(
		&user.ID,
		&user.Email,
		&user.RefreshToken,
		&user.PasswordHash,
		&user.RoleId,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	return &user, err
}
