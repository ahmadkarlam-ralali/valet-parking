package v1

import (
	"github.com/ahmadkarlam-ralali/valet-parking/helpers"
	"github.com/ahmadkarlam-ralali/valet-parking/models"
	"github.com/ahmadkarlam-ralali/valet-parking/requests"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"strconv"
)

type SlotsController struct {
	Db *gorm.DB
}

// List Slot by Building godoc
// @Summary List Slot by Building
// @Description list slot by building
// @Tags Slot
// @Accept  json
// @Produce  json
// @Param buildingID path string true "Building ID" default(1)
// @Success 200 {string} string "Ok"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /buildings/{buildingID}/slots [get]
func (this *SlotsController) GetAll(c *gin.Context) {
	var slots []models.Slot
	this.Db.Find(&slots, "building_id = ?", c.Param("buildingID"))
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

	var slot models.Slot
	if result := this.Db.First(&slot, "id = ?", c.Param("slotID")); result.Error != nil {
		helpers.HttpError(c, "Building not found", http.StatusNotFound)
		return
	}

	var count uint
	this.Db.Model(&models.Transaction{}).
		Where("slot_id = ? and end_at = '0000-00-00 00:00:00'", slot.ID).
		Count(&count)

	if request.Total < count {
		message := "Slot currently used"
		helpers.HttpError(c, message, http.StatusBadRequest)
		return
	}

	slot.Name = request.Name
	slot.Total = request.Total
	this.Db.Model(&slot).Updates(slot)

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
	this.Db.Delete(&models.Slot{}, "id = ? and building_id = ?", c.Param("slotID"), c.Param("buildingID"))

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
	type TotalData struct{ Total uint }
	var sumSlot TotalData
	var countTransaction TotalData
	this.Db.Table("slots").
		Select("sum(total) as total").
		Where("building_id = ?", c.Param("buildingID")).
		Scan(&sumSlot)
	this.Db.Table("slots").
		Select("count(plat_no) as total").
		Joins("left join transactions on transactions.slot_id = slots.id").
		Where("building_id = ? and end_at = '0000-00-00 00:00:00'", c.Param("buildingID")).
		Scan(&countTransaction)

	count := sumSlot.Total - countTransaction.Total
	message := "Current empty slot"
	if count < 1 {
		message = "Parking full"
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": message,
		"data": gin.H{
			"left": count,
		},
	})
}
