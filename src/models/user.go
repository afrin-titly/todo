package models

type User struct {
	ID       int    `json:"id,omitempty"`
	UserName string `json:"username"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password,omitempty" validate:"required,min=8"`
}
