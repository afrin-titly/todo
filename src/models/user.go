package models

type User struct {
	UserName string `json:"username"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password,omitempty" validate:"min=8"`
}
