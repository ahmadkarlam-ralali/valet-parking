package v1

import (
	"github.com/ahmadkarlam-ralali/valet-parking/helpers"
	"github.com/ahmadkarlam-ralali/valet-parking/models"
	"github.com/ahmadkarlam-ralali/valet-parking/repository"
	"github.com/ahmadkarlam-ralali/valet-parking/requests"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"strconv"
)

type SlotsController struct {
	Db             *gorm.DB
	SlotRepository repository.SlotRepository
}

// List Slot by Building godoc
// @Summary List Slot by Building
// @Description list slot by building
// @Tags Slot
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param buildingID path string true "Building ID" default(1)
// @Success 200 {string} string "Ok"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /buildings/{buildingID}/slots [get]
func (this *SlotsController) GetAll(c *gin.Context) {
	buildingID, _ := strconv.Atoi(c.Param("buildingID"))
	slots := this.SlotRepository.All(uint(buildingID))
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   slots,
	})
}

// Create Slot by Building godoc
// @Summary Create Slot by Building
// @Description create slot by building
// @Tags Slot
// @Accept  json
// @Produce  json
// @Param request body requests.SlotStoreRequest true "Request Body"
// @Param buildingID path string true "Building ID" default(1)
// @Success 200 {string} string "Ok"
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /buildings/{buildingID}/slots [post]
func (this *SlotsController) Store(c *gin.Context) {
	var request requests.SlotStoreRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		helpers.HttpError(c, "Wrong Format", http.StatusBadRequest)
		return
	}

	buildingId, _ := strconv.Atoi(c.Param("buildingID"))
	if result := this.Db.First(&models.Building{}, "id = ?", uint(buildingId)); result.Error != nil {
		helpers.HttpError(c, "Building not found", http.StatusNotFound)
		return
	}

	this.Db.Create(&models.Slot{
		Name:       request.Name,
		BuildingID: uint(buildingId),
		Total:      request.Total,
	})

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Slot created",
	})
}

// Update Slot by Building godoc
// @Summary Update Slot by Building
// @Description update slot by building
// @Tags Slot
// @Accept  json
// @Produce  json
// @Param request body requests.SlotUpdateRequest true "Request Body"
// @Param buildingID path string true "Building ID" default(1)
// @Param slotID path string true "Slot ID" default(1)
// @Success 200 {string} string "Ok"
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /buildings/{buildingID}/slots/{slotID} [put]
func (this *SlotsController) Update(c *gin.Context) {
	var request requests.SlotUpdateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		helpers.HttpError(c, "Wrong Format", http.StatusBadRequest)
		return
	}

	slotID, _ := strconv.Atoi(c.Param("slotID"))
	slot, err := this.SlotRepository.Get(uint(slotID))
	if err != nil {
		helpers.HttpError(c, "Building not found", http.StatusNotFound)
		return
	}

	count := this.SlotRepository.GetTotalSlotOccupied(uint(slotID))
	if request.Total < count {
		message := "Slot currently used"
		helpers.HttpError(c, message, http.StatusBadRequest)
		return
	}

	this.SlotRepository.Update(slot, request)

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Slot updated",
	})
}

// Delete Slot by Building godoc
// @Summary Delete Slot by Building
// @Description delete slot by building
// @Tags Slot
// @Accept  json
// @Produce  json
// @Param buildingID path string true "Building ID" default(1)
// @Param slotID path string true "Slot ID" default(1)
// @Success 200 {string} string "Ok"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /buildings/{buildingID}/slots/{slotID} [delete]
func (this *SlotsController) Destroy(c *gin.Context) {
	slotID, _ := strconv.Atoi(c.Param("slotID"))
	buildingID, _ := strconv.Atoi(c.Param("buildingID"))
	this.SlotRepository.DeleteByBuildingId(uint(slotID), uint(buildingID))

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Slot deleted",
	})
}

// Check Slot by Building godoc
// @Summary Check Slot by Building
// @Description check slot by building
// @Tags Slot
// @Accept  json
// @Produce  json
// @Param buildingID path string true "Building ID" default(1)
// @Success 200 {string} string "Ok"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /buildings/{buildingID}/slots/check [get]
func (this *SlotsController) Check(c *gin.Context) {
	buildingID, _ := strconv.Atoi(c.Param("buildingID"))
	remaining := this.SlotRepository.GetRemainingSlotByBuildingId(uint(buildingID))

	message := "Current empty slot"
	if remaining < 1 {
		message = "Parking full"
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": message,
		"data": gin.H{
			"remaining": remaining,
		},
	})
}
