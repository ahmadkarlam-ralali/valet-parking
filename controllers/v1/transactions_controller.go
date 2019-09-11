package v1

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/ahmadkarlam-ralali/valet-parking/helpers"
	"github.com/ahmadkarlam-ralali/valet-parking/models"
	"github.com/ahmadkarlam-ralali/valet-parking/requests"
	"github.com/ahmadkarlam-ralali/valet-parking/responses"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"math"
	"net/http"
	"time"
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

	h := md5.New()
	h.Write([]byte(fmt.Sprintf("%s%s%d", request.PlatNo, time.Now(), slot.ID)))
	code := hex.EncodeToString(h.Sum(nil))
	this.Db.Create(&models.Transaction{
		PlatNo:  request.PlatNo,
		SlotId:  slot.ID,
		StartAt: time.Now(),
		Code:    code,
	})

	slot.Status = "occupied"
	this.Db.Model(&slot).Updates(slot)

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": &responses.TransactionStartResponse{
			PlatNo:   request.PlatNo,
			SlotName: slot.Name,
			Code:     code,
		},
	})
}

func (this *TransactionsController) End(c *gin.Context) {
	var request requests.TransactionEndRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		helpers.HttpError(c, "Wrong Format", http.StatusBadRequest)
		return
	}

	var transaction models.Transaction
	if result := this.Db.First(&transaction, "code = ?", request.Code); result.Error != nil {
		helpers.HttpError(c, "Transaction not found", http.StatusNotFound)
		return
	}

	var slot models.Slot
	this.Db.First(&slot, "id = ?", transaction.SlotId)
	this.Db.Model(&slot).Update("status", "empty")

	transaction.EndAt = time.Now()
	duration := transaction.EndAt.Sub(transaction.StartAt).Hours()
	transaction.Total = uint(1500 * math.Ceil(duration))
	this.Db.Model(&transaction).Updates(transaction)

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   &responses.TransactionEndResponse{Total: transaction.Total},
	})
}
