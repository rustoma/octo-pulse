package validator

import (
	"github.com/go-playground/validator/v10"
	"github.com/rustoma/octo-pulse/internal/errors"
	"github.com/rustoma/octo-pulse/internal/models"
)

type authorValidator struct {
	validate *validator.Validate
}

func newAuthorValidator(validate *validator.Validate) *authorValidator {
	return &authorValidator{
		validate: validate,
	}
}

func (v *authorValidator) Validate(author *models.Author) error {
	err := v.validate.Struct(author)
	if err != nil {
		return errors.BadRequest{Err: err.Error()}
	}

	return nil
}
