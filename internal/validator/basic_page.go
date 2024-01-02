package validator

import (
	"github.com/go-playground/validator/v10"
	"github.com/rustoma/octo-pulse/internal/errors"
	"github.com/rustoma/octo-pulse/internal/models"
)

type basicPageValidator struct {
	validate *validator.Validate
}

func newBasicPageValidator(validate *validator.Validate) *basicPageValidator {
	return &basicPageValidator{
		validate: validate,
	}
}

func (v *basicPageValidator) Validate(category *models.BasicPage) error {
	err := v.validate.Struct(category)
	if err != nil {
		return errors.BadRequest{Err: err.Error()}
	}

	return nil
}
