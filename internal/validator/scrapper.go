package validator

import (
	"github.com/go-playground/validator/v10"
	"github.com/rustoma/octo-pulse/internal/errors"
	"github.com/rustoma/octo-pulse/internal/models"
)

type scrapperValidator struct {
	validate *validator.Validate
}

func newScrapperValidator(validate *validator.Validate) *scrapperValidator {
	return &scrapperValidator{
		validate: validate,
	}
}

type ScrapperValidation struct {
	Question    string `validate:"required"`
	Answear     string `validate:"required"`
	Href        string `validate:"required"`
	PageContent string `validate:"required"`
}

func (v *scrapperValidator) Validate(question *models.Question) error {
	propertiesToValidate := ScrapperValidation{
		Question:    question.Question,
		Answear:     question.Answer,
		Href:        question.Href,
		PageContent: question.PageContent,
	}

	err := v.validate.Struct(propertiesToValidate)

	if err != nil {
		return errors.BadRequest{Err: err.Error()}
	}

	return nil
}
