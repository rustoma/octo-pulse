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

func (s *SqlScrapperStore) GetQuestionSources(id int) ([]*models.QuestionSource, error) {

	stmt, args, err := sqlQb().
		Select("id_question_source, id_question, href, COALESCE(page_content, '') AS page_content").
		From("octopulse_question_sources").
		Where(squirrel.Eq{"id_question": id}).
		ToSql()

	if err != nil {
		logger.Err(err).Send()
		return nil, err
	}

	rows, err := s.DB.Query(stmt, args...)
	defer rows.Close()
	if err != nil {
		logger.Err(err).Send()
		return nil, err
	}

	var questionSources []*models.QuestionSource

	for rows.Next() {
		questionSource, err := scanToQuestionSource(rows)
		if err != nil {
			logger.Err(err).Send()
			return nil, err
		}

		questionSources = append(questionSources, questionSource)
	}

	return questionSources, nil
}

func (s *SqlScrapperStore) GetQuestion(id int) (*models.Question, error) {
	stmt, args, err := sqlQb().
		Select("id_question, question,COALESCE(answer, '') AS answer, href, COALESCE(page_content, '') AS page_content, octopulse_questions.fetched, id_category").
		From("octopulse_questions").
		Join("octopulse_phrases USING (id_phrase)").
		Where(squirrel.Eq{"id_question": id}).
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

	if question == nil {
		return nil, nil
	}

	questionSources, err := s.GetQuestionSources(id)

	question.Sources = questionSources

	return question, nil
}

func (s *SqlScrapperStore) GetQuestions(filters ...*storage.GetQuestionsFilters) ([]*models.Question, error) {

	questionsStatement := sqlQb().
		Select("id_question, question,COALESCE(answer, '') AS answer, href, COALESCE(page_content, '') AS page_content, octopulse_questions.fetched, id_category").
		From("octopulse_questions").
		Join("octopulse_phrases USING (id_phrase)").
		Limit(100)

	if len(filters) > 0 && filters[0].CategoryId != 0 {
		questionsStatement = questionsStatement.Where(squirrel.Eq{"id_category": filters[0].CategoryId, "octopulse_questions.fetched": 0})
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
		Update("octopulse_questions").
		SetMap(questionMap).
		Where(squirrel.Eq{"id_question": id}).
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
		&question.Answer,
		&question.Href,
		&question.PageContent,
		&question.Fetched,
		&question.CategoryId,
	)

	return &question, err
}

func scanToQuestionSource(rows *sql.Rows) (*models.QuestionSource, error) {
	var questionSource models.QuestionSource
	err := rows.Scan(
		&questionSource.Id,
		&questionSource.QuestionId,
		&questionSource.Href,
		&questionSource.PageContent,
	)

	return &questionSource, err
}

func convertQuestionToQuestionMap(question *models.Question) map[string]interface{} {
	return map[string]interface{}{
		"id_question":  question.Id,
		"question":     question.Question,
		"answer":       question.Answer,
		"href":         question.Href,
		"page_content": question.PageContent,
		"fetched":      question.Fetched,
	}
}
