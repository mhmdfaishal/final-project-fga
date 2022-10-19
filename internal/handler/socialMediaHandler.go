package handler

import (
	"final-project-fga/internal/domain"
	"final-project-fga/internal/helpers"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type SocialMediaHandler struct {
	db *gorm.DB
}

func NewSocialMediaHandler(db *gorm.DB) *SocialMediaHandler {
	return &SocialMediaHandler{
		db: db,
	}
}

func (s *SocialMediaHandler) Create(c *gin.Context) {
	db := s.db
	user := c.MustGet("user").(jwt.MapClaims)
	contentType := helpers.GetContentType(c)

	_, _ = db, contentType

	socialMedia := domain.SocialMedia{}
	payload := domain.CreateSocialMediaPayload{}
	if contentType == appJSON {
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid Argument",
				"error":   err.Error(),
			})
			return
		}
		validate := payload.Validate()
		if validate != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid Argument",
				"error":   validate,
			})
			return
		}
		socialMedia.Name = payload.Name
		socialMedia.SocialMediaURL = payload.SocialMediaURL
		socialMedia.UserID = int(user["id"].(float64))
	} else {
		c.ShouldBind(&socialMedia)
	}

	err := db.Create(&socialMedia).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"status": 201,
		"data": gin.H{
			"id":               socialMedia.ID,
			"name":             socialMedia.Name,
			"social_media_url": socialMedia.SocialMediaURL,
			"user_id":          socialMedia.UserID,
			"created_at":       socialMedia.CreatedAt,
		},
	})
}

func (s *SocialMediaHandler) GetSocialMedias(c *gin.Context) {
	db := s.db

	var socialMedias []domain.SocialMedia

	err := db.Find(&socialMedias).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	for i := range socialMedias {
		user := domain.User{}
		user, _ = s.GetUserById(socialMedias[i].UserID)
		photoUrl, _ := s.GetProfileImageUrl(socialMedias[i].UserID)
		userSocialMedia := domain.UserSocialMedia{
			ID:              socialMedias[i].UserID,
			Username:        user.Username,
			ProfileImageUrl: photoUrl,
		}
		socialMedias[i].User = &userSocialMedia
	}

	c.JSON(200, gin.H{
		"status": 200,
		"data": gin.H{
			"social_medias": socialMedias,
		},
	})
}

func (s *SocialMediaHandler) GetUserById(id int) (domain.User, error) {
	db := s.db

	var user domain.User

	err := db.Where("id = ?", id).First(&user).Error

	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (s *SocialMediaHandler) GetProfileImageUrl(id int) (string, error) {
	db := s.db

	var photo domain.Photo

	err := db.Where("user_id = ?", id).Order("created_at asc").First(&photo).Error

	if err != nil {
		return "", err
	}

	return photo.PhotoURL, nil
}

func (s *SocialMediaHandler) Update(c *gin.Context) {
	db := s.db

	contentType := helpers.GetContentType(c)

	_, _ = db, contentType

	socialMedia := domain.SocialMedia{}
	id := c.Param("socialMediaId")

	payload := domain.UpdateSocialMediaPayload{}
	if contentType == appJSON {
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid Argument",
				"error":   err.Error(),
			})
			return
		}

		validate := payload.Validate()
		if validate != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid Argument",
				"error":   validate,
			})
			return
		}

		socialMedia.Name = payload.Name
		socialMedia.SocialMediaURL = payload.SocialMediaURL
	} else {
		c.ShouldBind(&socialMedia)
	}

	err := db.Model(&socialMedia).Where("id = ?", id).Updates(&socialMedia).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = db.Where("id = ?", id).First(&socialMedia).Error
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"status": 200,
		"data": gin.H{
			"id":               socialMedia.ID,
			"name":             socialMedia.Name,
			"social_media_url": socialMedia.SocialMediaURL,
			"user_id":          socialMedia.UserID,
			"created_at":       socialMedia.CreatedAt,
		},
	})
}

func (s *SocialMediaHandler) Delete(c *gin.Context) {
	db := s.db

	id := c.Param("socialMediaId")

	socialMedia := domain.SocialMedia{}

	err := db.Where("id = ?", id).Delete(&socialMedia).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"status": 200,
		"data": gin.H{
			"message": "Your social media has been successfully deleted",
		},
	})
}
