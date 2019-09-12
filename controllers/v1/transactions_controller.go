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
	"log"
	"math"
	"net/http"
	"time"
)

type TransactionsController struct {
	Db *gorm.DB
}

// StartParking godoc
// @Summary Start Parking
// @Description start parking
// @Tags Transaction
// @Accept  json
// @Produce  json
// @Param request body requests.TransactionStartRequest true "Request Body"
// @Success 200 {string} string "Ok"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /transactions/start [post]
func (this *TransactionsController) Start(c *gin.Context) {
	var request requests.TransactionStartRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		helpers.HttpError(c, "Wrong Format", http.StatusBadRequest)
		return
	}

	type TotalData struct{ Total int }
	var sumSlot TotalData
	var countTransaction TotalData
	this.Db.Table("slots").
		Select("sum(total) as total").
		Where("building_id = ?", request.BuildingID).
		Scan(&sumSlot)
	this.Db.Table("slots").
		Select("count(plat_no) as total").
		Joins("left join transactions on transactions.slot_id = slots.id").
		Where("building_id = ? and end_at = '0000-00-00 00:00:00'", request.BuildingID).
		Scan(&countTransaction)

	count := sumSlot.Total - countTransaction.Total
	if count < 1 {
		helpers.HttpError(c, "Parking full", http.StatusBadRequest)
		return
	}

	type Parking struct {
		SlotID         uint
		SlotName       string
		TotalParking   int
		TotalAvailable int
	}
	rows, _ := this.Db.Raw("select "+
		"s.id as slot_id, s.name as slot_name, "+
		"(select count(*) from transactions t where end_at = '0000-00-00 00:00:00' and t.slot_id = s.id) as total_parking, "+
		"s.total as total_available from slots s where building_id = ?", request.BuildingID).
		Rows()

	var parking Parking
	for rows.Next() {
		err := rows.Scan(&parking.SlotID, &parking.SlotName, &parking.TotalParking, &parking.TotalAvailable)
		if err != nil {
			log.Println(err)
			helpers.HttpError(c, "Something when wrong", http.StatusInternalServerError)
			return
		}

		if parking.TotalParking < parking.TotalAvailable {
			break
		}
	}

	h := md5.New()
	h.Write([]byte(fmt.Sprintf("%s%s%d", request.PlatNo, time.Now(), parking.SlotID)))
	code := hex.EncodeToString(h.Sum(nil))
	this.Db.Create(&models.Transaction{
		PlatNo:  request.PlatNo,
		SlotId:  parking.SlotID,
		StartAt: time.Now(),
		Code:    code,
	})

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": &responses.TransactionStartResponse{
			PlatNo:   request.PlatNo,
			SlotName: parking.SlotName,
			Code:     code,
		},
	})
}

// EndParking godoc
// @Summary End Parking
// @Description end parking
// @Tags Transaction
// @Accept  json
// @Produce  json
// @Param request body requests.TransactionEndRequest true "Request Body"
// @Success 200 {string} string "Ok"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /transactions/end [post]
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

	transaction.EndAt = time.Now()
	duration := transaction.EndAt.Sub(transaction.StartAt).Hours()
	transaction.Total = int(1500 * math.Ceil(duration))
	this.Db.Model(&transaction).Updates(transaction)

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   &responses.TransactionEndResponse{Total: transaction.Total},
	})
}
