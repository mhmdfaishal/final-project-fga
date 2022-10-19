package middlewares

import (
	"errors"
	"final-project-fga/internal/domain"
	"fmt"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func Authorization(param, domain string, validate func(id, userID uint) (int, error)) gin.HandlerFunc {
	
	return func(c *gin.Context) {
		id := c.Param(param + "Id")
		var parseId, err = strconv.ParseUint(id, 10, 32)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "Invalid type of " + domain + " id",
			})
			return
		}

		user := c.MustGet("user").(jwt.MapClaims)  
		userID := uint(user["id"].(float64))

		if code, err := validate(uint(parseId), userID); err != nil {
			c.AbortWithStatusJSON(code, gin.H{
				"message": err.Error(),
			})
			return
		}

		c.Next()
	}
}

func CommentAuthorization(db *gorm.DB) gin.HandlerFunc {
	
	checkUserComment := func(id, userID uint) (int, error) {
		var comment domain.Comment
		err := db.Select("user_id").First(&comment, id).Error

		if err != nil {
			return http.StatusBadRequest, fmt.Errorf(fmt.Sprintf("the %s id %d was not found", "comment", id))
		}

		if uint(comment.UserID) != userID {
			return http.StatusUnauthorized, errors.New("unauthorized")
		}

		return http.StatusOK, nil
	}

	return Authorization("comment", "comment", checkUserComment)
}

func PhotoAuthorization(db *gorm.DB) gin.HandlerFunc {
	
	checkUserPhoto := func(id uint, userID uint) (int, error) {
		var photo domain.Photo
		err := db.Select("user_id").First(&photo, id).Error

		if err != nil {
			return http.StatusBadRequest, fmt.Errorf(fmt.Sprintf("the %s id %d was not found", "photo", id))
		}

		if uint(photo.UserID) != userID {
			return http.StatusUnauthorized, errors.New("unauthorized")
		}

		return http.StatusOK, nil
	}

	return Authorization("photo", "photo", checkUserPhoto)
}

func SocialMediaAuthorization(db *gorm.DB) gin.HandlerFunc {

	checkUserSocialMedia := func(id, userID uint) (int, error) {
		var socialMedia domain.SocialMedia
		var err = db.Select("user_id").First(&socialMedia, id).Error

		if err != nil {
			return http.StatusBadRequest, fmt.Errorf(fmt.Sprintf("the %s id %d was not found", "social media", id))
		}

		if uint(socialMedia.UserID) != userID {
			return http.StatusUnauthorized, errors.New("unauthorized")
		}

		return http.StatusOK, nil
	}

	return Authorization("socialMedia", "social media", checkUserSocialMedia)
}