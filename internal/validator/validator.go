package validator

import (
	"github.com/go-playground/validator/v10"
	"github.com/rustoma/octo-pulse/internal/models"
)

type Validator struct {
	Article ArticleValidatorer
}

func NewValidator() *Validator {
	validate := validator.New(validator.WithRequiredStructEnabled())

	return &Validator{
		Article: newArticleValidator(validate),
	}
}

type ArticleValidatorer interface {
	Validate(article *models.Article) error
}
