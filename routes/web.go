package routes

import (
	v1Controller "github.com/ahmadkarlam-ralali/valet-parking/controllers/v1"
	"github.com/ahmadkarlam-ralali/valet-parking/middlewares"
	"github.com/ahmadkarlam-ralali/valet-parking/repository"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.Use(middlewares.CORSMiddleware())

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	v1 := r.Group("/api/v1")

	buildingController := &v1Controller.BuildingsController{
		BuildingRepository: repository.BuildingRepository{Db: db},
	}
	slot := &v1Controller.SlotsController{
		BuildingRepository: repository.BuildingRepository{Db: db},
		SlotRepository:     repository.SlotRepository{Db: db},
	}
	transaction := &v1Controller.TransactionsController{
		SlotRepository:    repository.SlotRepository{Db: db},
		ParkingRepository: repository.ParkingRepository{Db: db},
	}
	report := &v1Controller.ReportsController{
		ParkingRepository: repository.ParkingRepository{Db: db},
	}

	authRoute := v1.Group("/auth")
	{
		auth := &v1Controller.AuthController{Db: db}
		authRoute.POST("/login", auth.Login)
	}

	// Without Authentication
	v1.GET("/buildings", buildingController.GetAll)
	v1.GET("/buildings/:buildingID/slots/check", slot.Check)
	v1.POST("/transactions/start", transaction.Start)
	v1.POST("/transactions/end", transaction.End)

	// With Authentication
	v1.Use(middlewares.Authenticate(db))
	{
		v1.POST("/buildings", buildingController.Store)
		v1.GET("/buildings/:buildingID", buildingController.Show)
		v1.PUT("/buildings/:buildingID", buildingController.Update)
		v1.DELETE("/buildings/:buildingID", buildingController.Destroy)

		v1.GET("/buildings/:buildingID/slots", slot.GetAll)
		v1.POST("/buildings/:buildingID/slots", slot.Store)
		v1.PUT("/buildings/:buildingID/slots/:slotID", slot.Update)
		v1.DELETE("/buildings/:buildingID/slots/:slotID", slot.Destroy)

		v1.GET("/transactions", transaction.GetAll)

		v1.GET("/reports/total-parking", report.GetTotalParking)
		v1.GET("/reports/total-income", report.GetIncomeParking)
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
