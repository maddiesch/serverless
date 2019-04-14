package serverless

import (
	"sync"

	"gopkg.in/go-playground/validator.v9"
)

var (
	validatorSetupNonce sync.Once
	validatorInstance   *validator.Validate
)

// GetValidator returns the default validator instance.
func GetValidator() *validator.Validate {
	validatorSetupNonce.Do(func() {
		validatorInstance = validator.New()
	})
	return validatorInstance
}

// Validatable is used to perform validation before a record is saved/updated
type Validatable interface {
	Validate() error
}
