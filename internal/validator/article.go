package validator

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/rustoma/octo-pulse/internal/errors"
	"github.com/rustoma/octo-pulse/internal/models"
)

type articleValidator struct {
	validate *validator.Validate
}

func newArticleValidator(validate *validator.Validate) *articleValidator {
	return &articleValidator{
		validate: validate,
	}
}

type ArticleValidation struct {
	Title           string `validate:"required,min=4" json:"title"`
	Description     string
	Thumbnail       *int
	PublicationDate time.Time
	IsPublished     bool
	AuthorId        int       `validate:"required"`
	CategoryId      int       `validate:"required"`
	DomainId        int       `validate:"required"`
	CreatedAt       time.Time `validate:"required"`
	UpdatedAt       time.Time `validate:"required"`
}

func (v *articleValidator) Validate(article *models.Article) error {
	propertiesToValidate := ArticleValidation{
		Title:           article.Title,
		Description:     article.Description,
		Thumbnail:       article.Thumbnail,
		PublicationDate: article.PublicationDate,
		IsPublished:     article.IsPublished,
		AuthorId:        article.AuthorId,
		CategoryId:      article.CategoryId,
		DomainId:        article.DomainId,
		CreatedAt:       article.CreatedAt,
		UpdatedAt:       article.UpdatedAt,
	}

	err := v.validate.Struct(propertiesToValidate)

	if err != nil {
		return errors.BadRequest{Err: err.Error()}
	}

	return nil
}
