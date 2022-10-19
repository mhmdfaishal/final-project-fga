package handler

import (
	"final-project-fga/internal/domain"
	"final-project-fga/internal/helpers"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type AuthHandler struct {
	db *gorm.DB
}

var (
	appJSON = "application/json"
)

func NewAuthHandler(db *gorm.DB) *AuthHandler {
	return &AuthHandler{
		db: db,
	}
}

func (ah *AuthHandler) Login(c *gin.Context) {
	db := ah.db

	contentType := helpers.GetContentType(c)
	_, _ = db, contentType

	User := domain.User{}
	password := ""

	if contentType == appJSON {
		c.ShouldBindJSON(&User)
	} else {
		c.ShouldBind(&User)
	}
	password = User.Password

	err := db.Where("email = ?", User.Email).First(&User).Error

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "Invalid email or password",
		})
		return
	}

	err = helpers.ComparePassword(User.Password, password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "Invalid email or password",
		})
		return
	}

	token := helpers.GenerateToken(User.ID, User.Email)

	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"data": gin.H{
			"token": token,
		},
	})
}

func (ah *AuthHandler) Register(c *gin.Context) {
	db := ah.db

	contentType := helpers.GetContentType(c)
	_, _ = db, contentType

	User := domain.User{}
	payload := domain.RegisterPayload{}

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
		// Hash password
		hashedPassword, err := helpers.HashPassword([]byte(payload.Password))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid Argument",
				"message": err.Error(),
			})
			return
		}
		// Assign hashed password to User
		User.Password = string(hashedPassword)
		User.Email = payload.Email
		User.Username = payload.Username
		User.Age = payload.Age

	} else {
		c.ShouldBind(&User)
	}

	err := db.Create(&User).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status": 201,
		"data": gin.H{
			"age":      User.Age,
			"email":    User.Email,
			"id":       User.ID,
			"username": User.Username,
		},
	})
}
