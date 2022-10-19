package handler

import (
	"final-project-fga/internal/domain"
	"final-project-fga/internal/helpers"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type PhotoHandler struct {
	db *gorm.DB
}

func NewPhotoHandler(db *gorm.DB) *PhotoHandler {
	return &PhotoHandler{
		db: db,
	}
}

func (ph *PhotoHandler) Create(c *gin.Context) {
	db := ph.db
	user := c.MustGet("user").(jwt.MapClaims)
	contentType := helpers.GetContentType(c)

	_, _ = db, contentType

	Photo := domain.Photo{}
	payload := domain.CreatePhotoPayload{}

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

		Photo.Title = payload.Title
		Photo.Caption = payload.Caption
		Photo.PhotoURL = payload.PhotoURL
		Photo.UserID = int(user["id"].(float64))
	} else {
		c.ShouldBind(&Photo)
	}

	err := db.Create(&Photo).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"status": 201,
		"data": gin.H{
			"id":         Photo.ID,
			"title":      Photo.Title,
			"caption":    Photo.Caption,
			"photo_url":  Photo.PhotoURL,
			"user_id":    Photo.UserID,
			"created_at": Photo.CreatedAt,
		},
	})
}

func (ph *PhotoHandler) GetPhotos(c *gin.Context) {
	db := ph.db

	var Photos []domain.Photo

	err := db.Find(&Photos).Error

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	for i := range Photos {
		user := domain.User{}
		user, _ = ph.GetUserById(Photos[i].UserID)
		userPhoto := domain.UserPhotos{
			Email:    user.Email,
			Username: user.Username,
		}
		Photos[i].User = &userPhoto
	}

	c.JSON(200, gin.H{
		"status": 200,
		"data":   Photos,
	})
}

func (ph *PhotoHandler) GetUserById(id int) (domain.User, error) {
	db := ph.db

	var user domain.User

	err := db.Where("id = ?", id).First(&user).Error

	if err != nil {
		return user, err
	}

	return user, nil
}

func (ph *PhotoHandler) Update(c *gin.Context) {
	db := ph.db
	user := c.MustGet("user").(jwt.MapClaims)
	contentType := helpers.GetContentType(c)

	_, _ = db, contentType

	id := c.Param("photoId")

	Photo := domain.Photo{}
	payload := domain.UpdatePhotoPayload{}
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

		Photo.Title = payload.Title
		Photo.Caption = payload.Caption
		Photo.PhotoURL = payload.PhotoURL
		Photo.UserID = int(user["id"].(float64))
	} else {
		c.ShouldBind(&Photo)
	}

	err := db.Model(&Photo).Where("id = ?", id).Updates(&Photo).Error

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = db.Where("id = ?", id).First(&Photo).Error
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"status": 200,
		"data": gin.H{
			"id":         Photo.ID,
			"title":      Photo.Title,
			"caption":    Photo.Caption,
			"photo_url":  Photo.PhotoURL,
			"user_id":    Photo.UserID,
			"updated_at": Photo.UpdatedAt,
		},
	})
}

func (ph *PhotoHandler) Delete(c *gin.Context) {
	db := ph.db

	id := c.Param("photoId")

	err := db.Where("id = ?", id).Delete(&domain.Photo{}).Error

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"status": 200,
		"data": gin.H{
			"message": "Your photo has been successfully deleted",
		},
	})
}
