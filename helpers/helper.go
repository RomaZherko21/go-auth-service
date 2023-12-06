package helpers

import (
	"os"

	"github.com/go-playground/validator"
	log "github.com/sirupsen/logrus"

	"github.com/joho/godotenv"
)

func GetEnv(key string) string {
	err := godotenv.Load(".env")

	if err != nil {
		log.Panicf("Error loading .env file: %v", err)
	}

	return os.Getenv(key)
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ValidationResult struct {
	OK     bool              `json:"ok"`
	Errors []ValidationError `json:"errors,omitempty"`
}

func Validate(data interface{}) ValidationResult {
	validate := validator.New()
	err := validate.Struct(data)

	if err != nil {
		var validationErrors []ValidationError
		for _, e := range err.(validator.ValidationErrors) {
			err := ValidationError{
				Field:   e.Field(),
				Message: e.Tag(),
			}
			validationErrors = append(validationErrors, err)
		}

		return ValidationResult{
			OK:     false,
			Errors: validationErrors,
		}
	}

	return ValidationResult{
		OK: true,
	}
}
