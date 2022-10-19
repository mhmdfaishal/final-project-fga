package handler

import (
	"final-project-fga/internal/domain"
	"final-project-fga/internal/helpers"
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type UserHandler struct {
	db *gorm.DB
}

func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{
		db: db,
	}
}

func (u *UserHandler) Update(c *gin.Context) {
	db := u.db

	contentType := helpers.GetContentType(c)
	_, _ = db, contentType

	User := domain.User{}
	id := c.Param("userId")

	payload := domain.UpdateUserPayload{}
	if contentType == appJSON {
		if err := c.BindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid Argument",
				"message": err.Error(),
			})
			return
		}

		validate := payload.Validate()
		if validate != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid Argument",
				"message": validate,
			})
			return
		}
		log.Println("USER >> ", payload)
	} else {
		c.ShouldBind(&User)
	}

	err := db.Model(&User).Where("id = ?", id).Updates(&payload).Error

	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = db.Where("id = ?", id).First(&User).Error
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"data": gin.H{
			"id":         User.ID,
			"email":      User.Email,
			"username":   User.Username,
			"age":        User.Age,
			"updated_at": User.UpdatedAt,
		},
	})
}

func (u *UserHandler) Delete(c *gin.Context) {
	db := u.db

	user := c.MustGet("user").(jwt.MapClaims)

	err := db.Where("id = ?", user["id"]).Delete(&domain.User{}).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"data": gin.H{
			"message": "Your account has been successfully deleted",
		},
	})
}
