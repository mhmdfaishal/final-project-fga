package database

import (
	"fmt"
	"log"

	"final-project-fga/internal/domain"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/spf13/viper"
)

func NewMysqlClient() *gorm.DB {
	dbHost := viper.GetString("DB_HOST")
	dbUser := viper.GetString("DB_USER")
	dbPass := viper.GetString("DB_PASS")
	dbName := viper.GetString("DB_NAME")
	dbPort := viper.GetString("DB_PORT")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", dbUser, dbPass, dbHost, dbPort, dbName)
	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		log.Println("Failed to connect to database ")
		panic(err)
	}

	log.Println("Connected to database")
	db.AutoMigrate(domain.SocialMedia{}, domain.User{}, domain.Photo{}, domain.Comment{})

	return db
}
