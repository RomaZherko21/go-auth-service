package helpers

import (
	"fmt"
	"unicode"

	"github.com/go-playground/validator/v10"
	"github.com/ttacon/libphonenumber"
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ValidationResult struct {
	OK     bool              `json:"ok"`
	Errors []ValidationError `json:"errors,omitempty"`
}

var validationMessages = map[string]string{
	"email":    "Email must contain @ and . symbols, f.e. user1@gmail.com",
	"password": "Password must contain at least one uppercase, lowercase, punctation symbol and numeric",
}

func Validate(data interface{}) ValidationResult {
	validate := validator.New()
	_ = validate.RegisterValidation("password", PasswordValidator)
	_ = validate.RegisterValidation("phone", PhoneNumberValidator)

	err := validate.Struct(data)

	if err != nil {
		var validationErrors []ValidationError
		for _, e := range err.(validator.ValidationErrors) {
			message, ok := validationMessages[e.Tag()]
			if !ok {
				message = fmt.Sprintf("validation err: %v", e.Tag())
			}

			err := ValidationError{
				Field:   e.Field(),
				Message: message,
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

func PhoneNumberValidator(fl validator.FieldLevel) bool {
	value := fl.Field().String()

	// Try to extract the country code from the phone number
	parsedNumber, err := libphonenumber.Parse(value, "")
	fmt.Println(parsedNumber)
	if err != nil {
		return false
	}

	return libphonenumber.IsValidNumber(parsedNumber)
}

func PasswordValidator(fl validator.FieldLevel) bool {
	value := fl.Field().String()

	var hasLowercase, hasUppercase, hasDigit, hasSpecialChar bool

	for _, char := range value {
		switch {
		case unicode.IsLower(char):
			hasLowercase = true
		case unicode.IsUpper(char):
			hasUppercase = true
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsSymbol(char) || unicode.IsPunct(char):
			hasSpecialChar = true
		}
	}

	return hasLowercase && hasUppercase && hasDigit && hasSpecialChar
}
