package routes

import (
	"github.com/ahmadkarlam-ralali/valet-parking/models"
	"github.com/jinzhu/gorm"
	"log"
	"os"
)

type Cmd struct {
	Db *gorm.DB
}

func (cmd *Cmd) HandleCmd() {
	firstArg := os.Args[1]
	switch firstArg {
	case "migrate":
		cmd.migrate()
		break
	default:
		log.Println("Nothing to execute.")
	}
}

func (cmd *Cmd) migrate() {
	cmd.Db.AutoMigrate(&models.Slot{}, &models.User{}, &models.Transaction{}, &models.Building{})
	log.Println("Migrate success")
}
