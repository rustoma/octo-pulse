package sqlstore

import (
	"database/sql"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/rustoma/octo-pulse/internal/models"
	"github.com/rustoma/octo-pulse/internal/storage"
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
		Select("id_phrase_result, question,COALESCE(answer, '') AS answer, href, COALESCE(page_content, '') AS page_content, octopulse_phrase_results.fetched, id_category").
		From("octopulse_phrase_results").
		Join("octopulse_phrases USING (id_phrase)").
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

func (s *SqlScrapperStore) GetQuestions(filters ...*storage.GetQuestionsFilters) ([]*models.Question, error) {

	questionsStatement := sqlQb().
		Select("id_phrase_result, question,COALESCE(answer, '') AS answer, href, COALESCE(page_content, '') AS page_content, octopulse_phrase_results.fetched, id_category").
		From("octopulse_phrase_results").
		Join("octopulse_phrases USING (id_phrase)").
		Limit(100)

	if len(filters) > 0 && filters[0].CategoryId != 0 {
		questionsStatement = questionsStatement.Where(squirrel.Eq{"id_category": filters[0].CategoryId, "octopulse_phrase_results.fetched": 0})
	}

	stmt, args, err := questionsStatement.ToSql()

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
	var questions []*models.Question

	for rows.Next() {
		questionFromScan, err := scanToQuestion(rows)

		if err != nil {
			logger.Err(err).Send()
			return nil, err
		}

		questions = append(questions, questionFromScan)
	}

	return questions, nil
}

func (s *SqlScrapperStore) UpdateQuestion(id int, question *models.Question) error {

	questionMap := convertQuestionToQuestionMap(question)

	stmt, args, err := sqlQb().
		Update("octopulse_phrase_results").
		SetMap(questionMap).
		Where(squirrel.Eq{"id_phrase_result": id}).
		ToSql()

	if err != nil {
		return err
	}

	_, err = s.DB.Exec(stmt, args...)

	return err
}

func scanToQuestion(rows *sql.Rows) (*models.Question, error) {
	var question models.Question
	err := rows.Scan(
		&question.Id,
		&question.Question,
		&question.Answear,
		&question.Href,
		&question.PageContent,
		&question.Fetched,
		&question.CategoryId,
	)

	return &question, err
}

func convertQuestionToQuestionMap(question *models.Question) map[string]interface{} {
	return map[string]interface{}{
		"id_phrase_result": question.Id,
		"question":         question.Question,
		"answer":           question.Answear,
		"href":             question.Href,
		"page_content":     question.PageContent,
		"fetched":          question.Fetched,
	}
}
