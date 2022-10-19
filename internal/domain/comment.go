package domain

import (
	"final-project-fga/pkg/validation"
	"time"
)

type Comment struct {
	ID        int           `json:"id"`
	Message   string        `json:"message"`
	PhotoID   int           `json:"photo_id"`
	UserID    int           `json:"user_id"`
	CreatedAt time.Time     `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time     `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
	User      *UserComment  `json:"User"`
	Photo     *PhotoComment `json:"Photo"`
}

type UserComment struct {
	ID       uint   `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

type PhotoComment struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Caption  string `json:"caption"`
	PhotoURL string `json:"photo_url"`
	UserID   int    `json:"user_id"`
}

type CreateCommentPayload struct {
	Message string `json:"message" validate:"required"`
	PhotoID int    `json:"photo_id" validate:"required"`
}

type UpdateCommentPayload struct {
	Message string `json:"message" validate:"required"`
}

func (p *CreateCommentPayload) Validate() []*validation.ErrorResponse {
	return validation.Validate(p)
}

func (p *UpdateCommentPayload) Validate() []*validation.ErrorResponse {
	return validation.Validate(p)
}
