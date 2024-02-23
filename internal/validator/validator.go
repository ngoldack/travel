package validator

import (
	"sync"

	"github.com/go-playground/validator/v10"
)

var (
	validate *validator.Validate
	once     sync.Once
)

func setup() {
	validate = validator.New(validator.WithRequiredStructEnabled())
}

func Get() *validator.Validate {
	once.Do(setup)
	return validate
}
