package sqlstore

import (
	"database/sql"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/rustoma/octo-pulse/internal/models"
)

type SqlScrapperStore struct {
	DB        *sql.DB
	dbTimeout time.Duration
}

func NewScrapperStore(DB *sql.DB) *SqlScrapperStore {
	return &SqlScrapperStore{
		DB:        DB,
		dbTimeout: time.Second * 3,
	}
}

func (s *SqlScrapperStore) GetQuestion(id int) (*models.Question, error) {

	stmt, args, err := sqlQb().
		Select("id_phrase_result, question, answer, href, page_content_processed").
		From("octopulse_phrase_results").
		Where(squirrel.Eq{"id_phrase_result": id}).
		ToSql()

	if err != nil {
		logger.Err(err).Send()
		return nil, err
	}

	rows, err := s.DB.Query(stmt, args...)

	if err != nil {
		logger.Err(err).Send()
		return nil, err
	}

	defer rows.Close()
	var question *models.Question

	for rows.Next() {
		questionFromScan, err := scanToQuestion(rows)

		if err != nil {
			logger.Err(err).Send()
			return nil, err
		}

		question = questionFromScan
	}

	return question, nil
}

func scanToQuestion(rows *sql.Rows) (*models.Question, error) {
	var question models.Question
	err := rows.Scan(
		&question.Id,
		&question.Question,
		&question.Answear,
		&question.Href,
		&question.PageContent,
	)

	return &question, err
}
