package main

import (
	"fmt"
	"github.com/ahmadkarlam-ralali/valet-parking/routes"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"log"
	"os"
)

func init() {
	viper.SetConfigFile("config.json")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func main() {
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

	if len(os.Args[1:]) > 0 && os.Args[1] == "migrate" {
		log.Println("Migrate success")
	} else {
		r := routes.SetupRouter(db)

		port := viper.GetString("app.port")

		_ = r.Run(fmt.Sprintf(":%s", port))
	}

	defer db.Close()
}
