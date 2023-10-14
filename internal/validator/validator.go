package validator

import (
	"github.com/go-playground/validator/v10"
	"github.com/rustoma/octo-pulse/internal/models"
)

type Validator struct {
	Article  ArticleValidatorer
	Scrapper ScrapperValidatorer
}

func NewValidator() *Validator {
	validate := validator.New(validator.WithRequiredStructEnabled())

	return &Validator{
		Article:  newArticleValidator(validate),
		Scrapper: newScrapperValidator(validate),
	}
}

type ArticleValidatorer interface {
	Validate(article *models.Article) error
}

type ScrapperValidatorer interface {
	Validate(question *models.Question) error
}
