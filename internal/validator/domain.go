package validator

import (
	"github.com/go-playground/validator/v10"
	"github.com/rustoma/octo-pulse/internal/errors"
	"github.com/rustoma/octo-pulse/internal/models"
)

type domainValidator struct {
	validate *validator.Validate
}

func newDomainValidator(validate *validator.Validate) *domainValidator {
	return &domainValidator{
		validate: validate,
	}
}

func (v *domainValidator) Validate(domain *models.Domain) error {
	err := v.validate.Struct(domain)
	if err != nil {
		return errors.BadRequest{Err: err.Error()}
	}

	return nil
}
