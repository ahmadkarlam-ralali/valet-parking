package v1

import (
	"github.com/ahmadkarlam-ralali/valet-parking/helpers"
	"github.com/ahmadkarlam-ralali/valet-parking/repository"
	"github.com/ahmadkarlam-ralali/valet-parking/requests"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type BuildingsController struct {
	BuildingRepository repository.BuildingRepository
}

// List Buildings godoc
// @Summary List Buildings
// @Description list building
// @Tags Building
// @Accept  json
// @Produce  json
// @Success 200 {string} string "Ok"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /buildings [get]
func (this *BuildingsController) GetAll(c *gin.Context) {
	buildings := this.BuildingRepository.GetAll()
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   buildings,
	})
}

// Create Buildings godoc
// @Summary Create Buildings
// @Description create building
// @Tags Building
// @Accept  json
// @Produce  json
// @Param request body requests.BuildingStoreRequest true "Request Body"
// @Success 201 {string} string "Created"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /buildings [post]
func (this *BuildingsController) Store(c *gin.Context) {
	var request requests.BuildingStoreRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		helpers.HttpError(c, "Wrong Format", http.StatusBadRequest)
		return
	}

	this.BuildingRepository.Create(request)

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Building created",
	})
}

// Update Buildings godoc
// @Summary Update Buildings
// @Description update building
// @Tags Building
// @Accept  json
// @Produce  json
// @Param request body requests.BuildingUpdateRequest true "Request Body"
// @Param buildingID path string true "Building ID" default(1)
// @Success 200 {string} string "Ok"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /buildings/{buildingID} [put]
func (this *BuildingsController) Update(c *gin.Context) {
	var request requests.BuildingUpdateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		helpers.HttpError(c, "Wrong Format", http.StatusBadRequest)
		return
	}

	buildingId, _ := strconv.Atoi(c.Param("buildingID"))
	this.BuildingRepository.Update(uint(buildingId), request)

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Building updated",
	})
}

// Delete Buildings godoc
// @Summary Delete Buildings
// @Description delete building
// @Tags Building
// @Accept  json
// @Produce  json
// @Param buildingID path string true "Building ID" default(11)
// @Success 200 {string} string "Ok"
// @Failure 500 {string} string "Internal Server Error"
// @Router /buildings/{buildingID} [delete]
func (this *BuildingsController) Destroy(c *gin.Context) {
	buildingId, _ := strconv.Atoi(c.Param("buildingID"))
	this.BuildingRepository.Delete(uint(buildingId))

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Building deleted",
	})
}
