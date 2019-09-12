package routes

import (
	v1Controller "github.com/ahmadkarlam-ralali/valet-parking/controllers/v1"
	"github.com/ahmadkarlam-ralali/valet-parking/middlewares"
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

	buildingRoute := v1.Group("/buildings")
	{
		controller := &v1Controller.BuildingsController{Db: db}
		buildingRoute.GET("/", controller.GetAll)
		buildingRoute.POST("/", controller.Store)
		buildingRoute.PUT("/:id", controller.Update)
		buildingRoute.DELETE("/:id", controller.Destroy)
	}

	slotRoute := v1.Group("/slots")
	{
		slot := &v1Controller.SlotsController{Db: db}
		slotRoute.GET("/", slot.GetAll)
		slotRoute.POST("/", slot.Store)
		slotRoute.PUT("/:id", slot.Update)
		slotRoute.DELETE("/:id", slot.Destroy)

		slotRoute.GET("/check", slot.Check)
	}

	transactionRoute := v1.Group("/transactions")
	{
		transaction := &v1Controller.TransactionsController{Db: db}
		transactionRoute.POST("/start", transaction.Start)
		transactionRoute.POST("/end", transaction.End)
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
