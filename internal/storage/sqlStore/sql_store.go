package sqlstore

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/rs/zerolog"
	lr "github.com/rustoma/octo-pulse/internal/logger"
	"github.com/rustoma/octo-pulse/internal/storage"
)

var logger *zerolog.Logger

type SqlStore struct {
	Scrapper storage.ScrapperStore
}

func NewSqlStorage(DB *sql.DB) *SqlStore {
	return &SqlStore{
		Scrapper: NewScrapperStore(DB),
	}
}

func sqlQb() squirrel.StatementBuilderType {
	return squirrel.StatementBuilder
}

func init() {
	//Init logger
	l, logFile := lr.NewLogger()
	defer logFile.Close()
	logger = l
}
