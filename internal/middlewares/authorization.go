package middlewares

import (
	"errors"
	"fmt"
	"final-project-fga/internal/domain"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func Authorization(param, modelName string, validate func(id, userID uint) (int, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param(param + "Id")
		var parseId, err = strconv.ParseUint(id, 10, 32)

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "Invalid type of " + modelName + " id",
			})
			return
		}

		userData := ctx.MustGet("userData").(jwt.MapClaims)
		userID := uint(userData["id"].(float64))

		if code, err := validate(uint(parseId), userID); err != nil {
			ctx.AbortWithStatusJSON(code, gin.H{
				"message": err.Error(),
			})
			return
		}

		ctx.Next()
	}
}

func CommentAuthorization(db *gorm.DB) gin.HandlerFunc {
	var checkUserComment = func(id, userID uint) (int, error) {
		var comment domain.Comment
		err := db.Select("user_id").First(&comment, id).Error

		if err != nil {
			return http.StatusBadRequest,
				fmt.Errorf(fmt.Sprintf("the %s id %d was not found", "comment", id))
		}

		if uint(comment.UserID) != userID {
			return http.StatusUnauthorized,
				errors.New("you are not allowed")
		}

		return http.StatusOK, nil
	}

	return Authorization("comment", "comment", checkUserComment)
}

func PhotoAuthorization(db *gorm.DB) gin.HandlerFunc {
	var checkUserPhoto = func(id uint, userID uint) (int, error) {
		var photo domain.Photo
		err := db.Select("user_id").First(&photo, id).Error

		if err != nil {
			return http.StatusBadRequest,
				fmt.Errorf(fmt.Sprintf("the %s id %d was not found", "photo", id))
		}

		if uint(photo.UserID) != userID {
			return http.StatusUnauthorized,
				errors.New("you are not allowed")
		}

		return http.StatusOK, nil
	}

	return Authorization("photo", "photo", checkUserPhoto)
}

func SocialMediaAuthorization(db *gorm.DB) gin.HandlerFunc {

	var checkUserSocialMedia = func(id, userID uint) (int, error) {
		var socialMedia domain.SocialMedia
		var err = db.Select("user_id").First(&socialMedia, id).Error

		if err != nil {
			return http.StatusBadRequest,
				fmt.Errorf(fmt.Sprintf("the %s id %d was not found", "social media", id))
		}

		if uint(socialMedia.UserID) != userID {
			return http.StatusUnauthorized,
				errors.New("you are not allowed")
		}

		return http.StatusOK, nil
	}

	return Authorization("socialMedia", "social media", checkUserSocialMedia)
}