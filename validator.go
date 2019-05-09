package serverless

import (
	"regexp"
	"sync"

	"github.com/segmentio/ksuid"
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
		validatorInstance.RegisterValidation("ksuid", validateKsuid)
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

func validateKsuid(fl validator.FieldLevel) bool {
	_, err := ksuid.Parse(fl.Field().String())
	return err == nil
}
