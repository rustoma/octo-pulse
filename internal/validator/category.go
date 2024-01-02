package validator

import (
	"github.com/go-playground/validator/v10"
	"github.com/rustoma/octo-pulse/internal/errors"
	"github.com/rustoma/octo-pulse/internal/models"
)

type categoryValidator struct {
	validate *validator.Validate
}

func newCategoryValidator(validate *validator.Validate) *categoryValidator {
	return &categoryValidator{
		validate: validate,
	}
}

func (v *categoryValidator) Validate(category *models.Category) error {
	err := v.validate.Struct(category)
	if err != nil {
		return errors.BadRequest{Err: err.Error()}
	}

	return nil
}
