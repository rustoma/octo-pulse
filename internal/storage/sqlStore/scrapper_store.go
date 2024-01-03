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
		dbTimeout: time.Second * 20,
	}
}

func (s *SqlScrapperStore) GetQuestionSources(id int) ([]*models.QuestionSource, error) {

	stmt, args, err := sqlQb().
		Select("id_question_source, id_question, href").
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

func (s *SqlScrapperStore) GetQuestionPageContents(id int) ([]*models.QuestionPageContent, error) {

	stmt, args, err := sqlQb().
		Select("id_question_source, id_question, href, COALESCE(page_content, '') AS page_content, COALESCE(page_content_processed, '') AS page_content_processed").
		From("octopulse_page_contents").
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

	var questionPageContents []*models.QuestionPageContent

	for rows.Next() {
		questionPageContent, err := scanToQuestionPageContent(rows)
		if err != nil {
			logger.Err(err).Send()
			return nil, err
		}

		questionPageContents = append(questionPageContents, questionPageContent)
	}

	return questionPageContents, nil
}

func (s *SqlScrapperStore) GetQuestion(id int) (*models.Question, error) {
	stmt, args, err := sqlQb().
		Select("id_question, question,COALESCE(answer, '') AS answer, href, octopulse_questions.fetched, id_category").
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

	questionPageContents, err := s.GetQuestionPageContents(id)

	question.PageContents = questionPageContents

	return question, nil
}

func (s *SqlScrapperStore) GetQuestions(filters ...*storage.GetQuestionsFilters) ([]*models.Question, error) {

	questionsStatement := sqlQb().
		Select("id_question, question,COALESCE(answer, '') AS answer, href, octopulse_questions.fetched, id_category").
		From("octopulse_questions").
		Join("octopulse_phrases USING (id_phrase)").
		OrderBy("RAND()").
		Limit(100)

	if len(filters) > 0 && filters[0].CategoryId != 0 {
		questionsStatement = questionsStatement.Where(
			squirrel.And{
				squirrel.Eq{"id_category": filters[0].CategoryId},
				squirrel.Eq{"octopulse_questions.fetched": 0},
			})
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

func (s *SqlScrapperStore) GetQuestionCategories() ([]*models.QuestionCategory, error) {
	stmt, args, err := sqlQb().
		Select("*").
		From("octopulse_categories").
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

	categories := make([]*models.QuestionCategory, 0)

	for rows.Next() {
		categoryFromScan, err := scanToQuestionCategory(rows)

		if err != nil {
			logger.Err(err).Send()
			return nil, err
		}

		categories = append(categories, categoryFromScan)
	}

	return categories, nil
}

func scanToQuestion(rows *sql.Rows) (*models.Question, error) {
	var question models.Question
	err := rows.Scan(
		&question.Id,
		&question.Question,
		&question.Answer,
		&question.Href,
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
	)

	return &questionSource, err
}

func scanToQuestionPageContent(rows *sql.Rows) (*models.QuestionPageContent, error) {
	var questionPageContent models.QuestionPageContent

	err := rows.Scan(
		&questionPageContent.SourceId,
		&questionPageContent.QuestionId,
		&questionPageContent.Href,
		&questionPageContent.PageContent,
		&questionPageContent.PageContentProcessed,
	)

	return &questionPageContent, err
}

func scanToQuestionCategory(rows *sql.Rows) (*models.QuestionCategory, error) {
	var category models.QuestionCategory
	err := rows.Scan(
		&category.IdCategory,
		&category.Name,
	)

	return &category, err
}

func convertQuestionToQuestionMap(question *models.Question) map[string]interface{} {
	return map[string]interface{}{
		"id_question": question.Id,
		"question":    question.Question,
		"answer":      question.Answer,
		"href":        question.Href,
		"fetched":     question.Fetched,
	}
}
