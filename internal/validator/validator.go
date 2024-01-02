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
	Category      CategoryValidatorer
	BasicPage     BasicPageValidatorer
}

func NewValidator() *Validator {
	validate := validator.New(validator.WithRequiredStructEnabled())

	return &Validator{
		Article:       newArticleValidator(validate),
		Scrapper:      newScrapperValidator(validate),
		Domain:        newDomainValidator(validate),
		ImageCategory: newImageCategoryValidator(validate),
		Author:        newAuthorValidator(validate),
		Category:      newCategoryValidator(validate),
		BasicPage:     newBasicPageValidator(validate),
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
	Validate(author *models.Author) error
}

type CategoryValidatorer interface {
	Validate(category *models.Category) error
}

type BasicPageValidatorer interface {
	Validate(category *models.BasicPage) error
}
