package validations

import (
	"fmt"
	"strconv"
	"todo-list/src/models"

	"github.com/go-playground/validator/v10"
)

// var validateUser *validator.Validate

func init() {
	validate = validator.New()
}

func ValidateUser(user *models.User) map[string]string {
	errors := make(map[string]string)
	err := validate.Struct(user)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var errorMessage string

			switch err.Tag() {
			case "required":
				errorMessage = "This field is required"
			case "min":
				minValue, _ := strconv.Atoi(err.Param())
				errorMessage = fmt.Sprintf("This field must be longer than %d characters", minValue)
			case "email":
				errorMessage = "Not a valid email address"
			default:
				errorMessage = fmt.Sprintf("failed on the '%s' tag", err.Tag())
			}
			errors[err.Field()] = errorMessage
		}
	}
	return errors
}
