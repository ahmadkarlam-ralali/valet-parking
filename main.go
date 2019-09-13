package main

import (
	"fmt"
	"github.com/ahmadkarlam-ralali/valet-parking/docs"
	"github.com/ahmadkarlam-ralali/valet-parking/routes"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"log"
	"os"
)

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func init() {
	viper.SetConfigFile("config.json")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func main() {
	// programatically set swagger info
	docs.SwaggerInfo.Title = "Swagger Agent API"
	docs.SwaggerInfo.Description = "Valet Parking API"
	docs.SwaggerInfo.Version = "0.1"
	docs.SwaggerInfo.Host = ""
	docs.SwaggerInfo.BasePath = "/api/v1/"
	docs.SwaggerInfo.Host = "192.168.40.94:9000"

	dbHostname := viper.GetString("database.hostname")
	dbUsername := viper.GetString("database.username")
	dbPassword := viper.GetString("database.password")
	dbName := viper.GetString("database.name")
	dbPort := viper.GetString("database.port")
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		dbUsername, dbPassword, dbHostname, dbPort, dbName)
	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Connection to database failed: %s", err)
	}

	if len(os.Args[1:]) > 0 {
		cmd := routes.Cmd{Db: db}
		cmd.HandleCmd()
	} else {
		r := routes.SetupRouter(db)

		port := viper.GetString("app.port")

		_ = r.Run(fmt.Sprintf(":%s", port))
	}

	defer db.Close()
}
