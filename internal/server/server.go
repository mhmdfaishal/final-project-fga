package server

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"

	"final-project-fga/internal/handler"
	"final-project-fga/internal/middlewares"
)

type Server struct {
	contextTimeout time.Duration
}

func NewServer() *Server {
	context_timeout := viper.GetInt("CONTEXT_TIMEOUT")
	return &Server{
		contextTimeout: time.Duration(context_timeout) * time.Second,
	}
}

func (s *Server) Run(db *gorm.DB) {

	port := viper.GetString("PORT")
	route := s.Route(db)
	go func() {
		fmt.Printf("Server running on port: %s\n", port)
		route.Run(fmt.Sprintf(":%s", port))
	}()

	quit := make(chan os.Signal, 1) 	// Create a channel for os.Signal
	signal.Notify(quit, os.Interrupt) 	// Notify the channel when an interrupt is received
	<-quit 								// Block until an interrupt is received
	
	fmt.Println("shutting down...")		
}

func (s *Server) Route(db *gorm.DB) *gin.Engine {
	router := gin.Default()

	authHandler := handler.NewAuthHandler(db)
	userHandler := handler.NewUserHandler(db)

	userRouter := router.Group("/users")
	{
		userRouter.POST("/register", authHandler.Register)
		userRouter.POST("/login", authHandler.Login)

		userRouter.Use(middlewares.Authentication()).PUT("/:userId", userHandler.Update)
		userRouter.Use(middlewares.Authentication()).DELETE("/", userHandler.Delete)
	}

	photoHandler := handler.NewPhotoHandler(db)

	photoRouter := router.Group("/photos")
	{
		photoRouter.Use(middlewares.Authentication())
		photoRouter.POST("/", photoHandler.Create)
		photoRouter.GET("/", photoHandler.GetPhotos)
		photoRouter.Use(middlewares.PhotoAuthorization(db)).PUT("/:photoId", photoHandler.Update)
		photoRouter.Use(middlewares.PhotoAuthorization(db)).DELETE("/:photoId", photoHandler.Delete)
	}

	commentHandler := handler.NewCommentHandler(db)

	commentRouter := router.Group("/comments")
	{
		commentRouter.Use(middlewares.Authentication())
		commentRouter.POST("/", commentHandler.Create)
		commentRouter.GET("/", commentHandler.GetComments)
		commentRouter.Use(middlewares.CommentAuthorization(db)).PUT("/:commentId", commentHandler.Update)
		commentRouter.Use(middlewares.CommentAuthorization(db)).DELETE("/:commentId", commentHandler.Delete)
	}

	socialMediaHandler := handler.NewSocialMediaHandler(db)

	socialMediaRouter := router.Group("/socialmedias")
	{
		socialMediaRouter.Use(middlewares.Authentication())
		socialMediaRouter.POST("/", socialMediaHandler.Create)
		socialMediaRouter.GET("/", socialMediaHandler.GetSocialMedias)
		socialMediaRouter.Use(middlewares.SocialMediaAuthorization(db)).PUT("/:socialMediaId", socialMediaHandler.Update)
		socialMediaRouter.Use(middlewares.SocialMediaAuthorization(db)).DELETE("/:socialMediaId", socialMediaHandler.Delete)
	}

	return router
}
