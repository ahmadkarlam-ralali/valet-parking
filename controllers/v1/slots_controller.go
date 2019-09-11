package v1

import (
	"github.com/ahmadkarlam-ralali/valet-parking/models"
	"github.com/ahmadkarlam-ralali/valet-parking/requests"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
)

type SlotsController struct {
	Db *gorm.DB
}

func (this *SlotsController) GetAll(c *gin.Context) {
	var slots []models.Slot
	this.Db.Find(&slots)
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   slots,
	})
}

func (this *SlotsController) Store(c *gin.Context) {
	var request requests.SlotStoreRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.Abort()
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Wrong format",
		})
		return
	}

	this.Db.Create(&models.Slot{
		Name:      request.Name,
		Status:    "empty",
	})

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Slot created",
	})
}

func (this *SlotsController) Update(c *gin.Context) {
	var request requests.SlotUpdateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.Abort()
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Wrong format",
		})
		return
	}

	var slot models.Slot
	this.Db.First(&slot, "id = ?", c.Param("id"))
	if request.Name != "" {
		slot.Name = request.Name
	}
	if request.Status != "" {
		slot.Status = request.Status
	}
	this.Db.Model(&slot).Updates(slot)

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Slot updated",
	})
}

func (this *SlotsController) Destroy(c *gin.Context) {
	this.Db.Delete(&models.Slot{}, "id = ?", c.Param("id"))

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Slot deleted",
	})
}