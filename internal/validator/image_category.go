package validator

import (
	"github.com/go-playground/validator/v10"
	"github.com/rustoma/octo-pulse/internal/errors"
	"github.com/rustoma/octo-pulse/internal/models"
)

type imageCategoryValidator struct {
	validate *validator.Validate
}

func newImageCategoryValidator(validate *validator.Validate) *imageCategoryValidator {
	return &imageCategoryValidator{
		validate: validate,
	}
}

func (v *imageCategoryValidator) Validate(category *models.ImageCategory) error {
	err := v.validate.Struct(category)
	if err != nil {
		return errors.BadRequest{Err: err.Error()}
	}

	return nil
}
