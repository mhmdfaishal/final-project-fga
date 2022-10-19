package handler

import (
	"final-project-fga/internal/domain"
	"final-project-fga/internal/helpers"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type CommentHandler struct {
	db *gorm.DB
}

func NewCommentHandler(db *gorm.DB) *CommentHandler {
	return &CommentHandler{
		db: db,
	}
}

func (co *CommentHandler) Create(c *gin.Context) {
	db := co.db
	user := c.MustGet("user").(jwt.MapClaims)
	contentType := helpers.GetContentType(c)

	_, _ = db, contentType

	Comment := domain.Comment{}
	payload := domain.CreateCommentPayload{}

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

		Comment.Message = payload.Message
		Comment.PhotoID = payload.PhotoID
		Comment.UserID = int(user["id"].(float64))
	} else {
		c.ShouldBind(&Comment)
	}

	err := db.Create(&Comment).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"status": 201,
		"data": gin.H{
			"id":         Comment.ID,
			"message":    Comment.Message,
			"photo_id":   Comment.PhotoID,
			"user_id":    Comment.UserID,
			"created_at": Comment.CreatedAt,
		},
	})
}

func (co *CommentHandler) GetComments(c *gin.Context) {
	db := co.db

	var Comments []domain.Comment

	err := db.Find(&Comments).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	for i := range Comments {
		user := domain.User{}
		user, _ = co.GetUserById(Comments[i].UserID)

		userComment := domain.UserComment{
			ID:       user.ID,
			Email:    user.Email,
			Username: user.Username,
		}

		photo := domain.Photo{}
		photo, _ = co.GetPhotoByUserId(Comments[i].UserID)

		photoComment := domain.PhotoComment{
			ID:       photo.ID,
			Title:    photo.Title,
			Caption:  photo.Caption,
			PhotoURL: photo.PhotoURL,
			UserID:   photo.UserID,
		}

		Comments[i].User = &userComment
		Comments[i].Photo = &photoComment
	}

	c.JSON(200, gin.H{
		"status": 200,
		"data":   Comments,
	})
}

func (co *CommentHandler) GetUserById(id int) (domain.User, error) {
	db := co.db

	var user domain.User

	err := db.Where("id = ?", id).First(&user).Error

	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (co *CommentHandler) GetPhotoByUserId(id int) (domain.Photo, error) {
	db := co.db

	var photo domain.Photo

	err := db.Where("user_id = ?", id).First(&photo).Error

	if err != nil {
		return domain.Photo{}, err
	}

	return photo, nil
}

func (co *CommentHandler) Update(c *gin.Context) {
	db := co.db

	contentType := helpers.GetContentType(c)

	_, _ = db, contentType

	Comment := domain.Comment{}
	id := c.Param("commentId")
	payload := domain.UpdateCommentPayload{}

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
		Comment.Message = payload.Message
	} else {
		c.ShouldBind(&Comment)
	}

	err := db.Model(&Comment).Where("id = ?", id).Updates(Comment).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = db.Where("id = ?", id).First(&Comment).Error
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	Photo := domain.Photo{}
	err = db.Where("id = ?", Comment.PhotoID).First(&Photo).Error
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

func (co *CommentHandler) Delete(c *gin.Context) {
	db := co.db

	id := c.Param("commentId")

	err := db.Where("id = ?", id).Delete(&domain.Comment{}).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"status": 200,
		"data": gin.H{
			"message": "Your comment has been successfully deleted",
		},
	})
}
