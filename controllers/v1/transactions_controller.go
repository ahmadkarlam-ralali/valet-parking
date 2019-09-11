package v1

import (
	"github.com/ahmadkarlam-ralali/valet-parking/helpers"
	"github.com/ahmadkarlam-ralali/valet-parking/models"
	"github.com/ahmadkarlam-ralali/valet-parking/requests"
	"github.com/ahmadkarlam-ralali/valet-parking/responses"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
)

type TransactionsController struct {
	Db *gorm.DB
}

func (this *TransactionsController) Start(c *gin.Context) {
	var request requests.TransactionStartRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		helpers.HttpError(c, "Wrong Format", http.StatusBadRequest)
		return
	}

	var slot models.Slot
	if result := this.Db.First(&slot, "status = 'empty'"); result.Error != nil {
		helpers.HttpError(c, "Parking full", http.StatusBadRequest)
		return
	}

	this.Db.Create(&models.Transaction{PlatNo: request.PlatNo, SlotId: slot.ID})

	slot.Status = "occupied"
	this.Db.Model(&slot).Updates(slot)

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": &responses.TransactionStartResponse{
			PlatNo:   request.PlatNo,
			SlotName: slot.Name,
			SlotID:   slot.ID,
		},
	})
}

func (this *TransactionsController) End(c *gin.Context) {
	var request requests.TransactionEndRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		helpers.HttpError(c, "Wrong Format", http.StatusBadRequest)
		return
	}

	var slot models.Slot
	if result := this.Db.First(&slot, "name = ?", request.SlotName); result.Error != nil {
		helpers.HttpError(c, "Slot Name not found", http.StatusNotFound)
		return
	}

	var transaction models.Transaction
	condition := "plat_no = ? and slot_id = ?"
	if result := this.Db.First(&transaction, condition, request.PlatNo, slot.ID); result.Error != nil {
		helpers.HttpError(c, "Transaction not found", http.StatusNotFound)
		return
	}

	slot.Status = "empty"
	this.Db.Model(&slot).Updates(slot)

	transaction.Total = 1500
	this.Db.Model(&transaction).Updates(transaction)

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   &responses.TransactionEndResponse{Total: transaction.Total},
	})
}
