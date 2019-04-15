package serverless

import (
	"regexp"
	"sync"

	"gopkg.in/go-playground/validator.v9"
)

var (
	validatorSetupNonce    sync.Once
	validatorInstance      *validator.Validate
	validateErrorCodeRegex = regexp.MustCompile(`\A(?:[^\d_])[\w]{3,}(?:[^_])\z`)
)

// GetValidator returns the default validator instance.
func GetValidator() *validator.Validate {
	validatorSetupNonce.Do(func() {
		validatorInstance = validator.New()
		validatorInstance.RegisterValidation("error_code", validateErrorCode)
	})
	return validatorInstance
}

// Validatable is used to perform validation before a record is saved/updated
type Validatable interface {
	Validate() error
}

func validateErrorCode(fl validator.FieldLevel) bool {
	value := fl.Field().String()

	return validateErrorCodeRegex.MatchString(value)
}
