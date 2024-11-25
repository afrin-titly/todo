package validations

import (
	"fmt"
	"strconv"
	"todo-list/src/models"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func ValidateTodo(todo *models.Todo) map[string]string {
	errors := make(map[string]string)
	err := validate.Struct(todo)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var errorMessage string

			switch err.Tag() {
			case "required":
				errorMessage = "This field is required"
			case "min":
				minValue, _ := strconv.Atoi(err.Param())
				errorMessage = fmt.Sprintf("This field must be longer than %d characters", minValue)
			default:
				errorMessage = fmt.Sprintf("failed on the '%s' tag", err.Tag())
			}
			errors[err.Field()] = errorMessage
		}
	}
	return errors
}
