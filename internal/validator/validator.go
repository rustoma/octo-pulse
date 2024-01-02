package validator

import (
	"github.com/go-playground/validator/v10"
	"github.com/rustoma/octo-pulse/internal/models"
)

type Validator struct {
	Article       ArticleValidatorer
	Scrapper      ScrapperValidatorer
	Domain        DomainValidatorer
	ImageCategory ImageCategoryValidatorer
	Author        AuthorValidatorer
}

func NewValidator() *Validator {
	validate := validator.New(validator.WithRequiredStructEnabled())

	return &Validator{
		Article:       newArticleValidator(validate),
		Scrapper:      newScrapperValidator(validate),
		Domain:        newDomainValidator(validate),
		ImageCategory: newImageCategoryValidator(validate),
		Author:        newAuthorValidator(validate),
	}
}

type ArticleValidatorer interface {
	Validate(article *models.Article) error
}

type ScrapperValidatorer interface {
	Validate(question *models.Question) error
}

type DomainValidatorer interface {
	Validate(domain *models.Domain) error
}

type ImageCategoryValidatorer interface {
	Validate(category *models.ImageCategory) error
}

type AuthorValidatorer interface {
	Validate(category *models.Author) error
}
