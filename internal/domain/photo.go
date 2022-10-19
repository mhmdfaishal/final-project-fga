package domain

import (
	"final-project-fga/pkg/validation"
	"time"
)

type Photo struct {
	ID        int         `json:"id"`
	Title     string      `json:"title"`
	Caption   string      `json:"caption"`
	PhotoURL  string      `json:"photo_url"`
	UserID    int         `json:"user_id"`
	CreatedAt time.Time   `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time   `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
	User      *UserPhotos `json:"User"`
}

type UserPhotos struct {
	Email    string `json:"email"`
	Username string `json:"username"`
}

type CreatePhotoPayload struct {
	Title    string `json:"title" validate:"required"`
	Caption  string `json:"caption"`
	PhotoURL string `json:"photo_url" validate:"required"`
}

type UpdatePhotoPayload struct {
	Title    string `json:"title" validate:"required"`
	Caption  string `json:"caption"`
	PhotoURL string `json:"photo_url" validate:"required"`
}

func (p *CreatePhotoPayload) Validate() []*validation.ErrorResponse {
	return validation.Validate(p)
}

func (p *UpdatePhotoPayload) Validate() []*validation.ErrorResponse {
	return validation.Validate(p)
}
