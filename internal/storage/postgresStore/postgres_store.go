package postgresstore

import (
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	lr "github.com/rustoma/octo-pulse/internal/logger"
	"github.com/rustoma/octo-pulse/internal/storage"
)

var logger *zerolog.Logger

func NewPostgresStorage(DB *pgxpool.Pool) *storage.Store {
	return &storage.Store{
		User:     NewUserStore(DB),
		Role:     NewRoleStore(DB),
		Domain:   NewDomainStore(DB),
		Category: NewCategoryStore(DB),
		Author:   NewAuthorStore(DB),
		Article:  NewArticleStore(DB),
	}
}

func pgQb() squirrel.StatementBuilderType {
	return squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
}

func init() {
	//Init logger
	l, logFile := lr.NewLogger()
	defer logFile.Close()
	logger = l
}
