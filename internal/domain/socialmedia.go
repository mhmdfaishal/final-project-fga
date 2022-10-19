package domain

import (
	"final-project-fga/pkg/validation"
	"time"
)

type SocialMedia struct {
	ID             int              `json:"id"`
	Name           string           `json:"name"`
	SocialMediaURL string           `json:"social_media_url"`
	UserID         int              `json:"user_id"`
	CreatedAt      time.Time        `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt      time.Time        `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
	User           *UserSocialMedia `json:"User"`
}

type UserSocialMedia struct {
	ID              int    `json:"id"`
	Username        string `json:"username"`
	ProfileImageUrl string `json:"profile_image_url"`
}

type CreateSocialMediaPayload struct {
	Name           string `json:"name" validate:"required"`
	SocialMediaURL string `json:"social_media_url" validate:"required"`
}

type UpdateSocialMediaPayload struct {
	Name           string `json:"name" validate:"required"`
	SocialMediaURL string `json:"social_media_url" validate:"required"`
}

func (p *CreateSocialMediaPayload) Validate() []*validation.ErrorResponse {
	return validation.Validate(p)
}

func (p *UpdateSocialMediaPayload) Validate() []*validation.ErrorResponse {
	return validation.Validate(p)
}
