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

	transactionRoute := v1.Group("/transactions")
	{
		transactionRoute.POST("/start", transaction.Start)
		transactionRoute.POST("/end", transaction.End)
	}

	buildingRoute := v1.Group("/buildings")
	buildingRoute.GET("/", buildingController.GetAll)

	slotRoute := buildingRoute.Group("/:buildingID/slots")
	slotRoute.GET("/check", slot.Check)

	// With Authentication
	v1.Use(middlewares.Authenticate(db))
	{
		buildingRoute.POST("/", buildingController.Store)
		buildingRoute.PUT("/:buildingID", buildingController.Update)
		buildingRoute.DELETE("/:buildingID", buildingController.Destroy)

		slotRoute.GET("/", slot.GetAll)
		slotRoute.POST("/", slot.Store)
		slotRoute.PUT("/:slotID", slot.Update)
		slotRoute.DELETE("/:slotID", slot.Destroy)

		transactionRoute.GET("/", transaction.GetAll)
	}

	authRoute := v1.Group("/auth")
	{
		auth := &v1Controller.AuthController{Db: db}
		authRoute.POST("/login", auth.Login)
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
