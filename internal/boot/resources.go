package boot

import (
	"fmt"
	"log"

	mysql "final-project-fga/pkg/database"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
)

var (
	DB *gorm.DB
)

func init() {
	viper.SetConfigFile(`.env`)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	if viper.GetBool("DEBUG") {
		fmt.Println("Service RUN on DEBUG mode")
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	DB = mysql.NewMysqlClient()
}

func FlushResources() {

	fmt.Println("stopping db connection")
	err := DB.Close()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("db connection stopped")
}
