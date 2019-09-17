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
		viper.SetDefault("database.hostname", "")
		viper.SetDefault("database.username", "")
		viper.SetDefault("database.password", "")
		viper.SetDefault("database.name", "")
		viper.SetDefault("database.port", "")
		viper.SetDefault("app.port", "9000")
		viper.SetDefault("app.env", "production")
		viper.SetDefault("app.env", "production")
		viper.SetDefault("swagger.host", "https://valet-parking-go.herokuapp.com")
	}
}

func main() {
	// programatically set swagger info
	docs.SwaggerInfo.Title = "Swagger Agent API"
	docs.SwaggerInfo.Description = "Valet Parking API"
	docs.SwaggerInfo.Version = "0.1"
	docs.SwaggerInfo.Host = ""
	docs.SwaggerInfo.BasePath = "/api/v1/"
	docs.SwaggerInfo.Host = viper.GetString("swagger.host")

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
		env := viper.GetString("app.env")

		if env == "production" {
			_ = r.Run(fmt.Sprintf(":%s", "80"))
		} else {
			_ = r.Run(fmt.Sprintf(":%s", port))
		}
	}

	defer db.Close()
}
