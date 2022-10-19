package domain

import (
	"final-project-fga/pkg/validation"
	"time"
)

type User struct {
	ID              uint      `json:"id"`
	Username        string    `json:"username" gorm:"unique"`
	Email           string    `json:"email" gorm:"unique"`
	Password        string    `json:"password"`
	Age             int       `json:"age"`
	CreatedAt       time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt       time.Time `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
}

type RegisterPayload struct {
	Age      int    `json:"age" validate:"required,min=9"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Username string `json:"username" validate:"required"`
}

type UpdateUserPayload struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
}

func (p *RegisterPayload) Validate() []*validation.ErrorResponse {
	return validation.Validate(p)
}

func (p *UpdateUserPayload) Validate() []*validation.ErrorResponse {
	return validation.Validate(p)
}
