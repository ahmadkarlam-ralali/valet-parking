package v1

import (
	"github.com/ahmadkarlam-ralali/valet-parking/helpers"
	"github.com/ahmadkarlam-ralali/valet-parking/repository"
	"github.com/ahmadkarlam-ralali/valet-parking/requests"
	"github.com/ahmadkarlam-ralali/valet-parking/responses"
	"github.com/gin-gonic/gin"
	"net/http"
)

type TransactionsController struct {
	SlotRepository    repository.SlotRepository
	ParkingRepository repository.ParkingRepository
}

// List Transaction godoc
// @Summary List Transaction
// @Description list transaction
// @Tags Transaction
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {string} string "Ok"
// @Failure 500 {string} string "Internal Server Error"
// @Router /transactions/ [get]
func (this *TransactionsController) GetAll(c *gin.Context) {
	transactions := this.ParkingRepository.GetAll()
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   transactions,
	})
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
// @Router /transactions/start/ [post]
func (this *TransactionsController) Start(c *gin.Context) {
	var request requests.TransactionStartRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		helpers.HttpError(c, "Wrong Format", http.StatusBadRequest)
		return
	}

	isStillParking := this.ParkingRepository.IsPlatNoStillParking(request.PlatNo)
	if isStillParking {
		helpers.HttpError(c, "This license number still parking", http.StatusBadRequest)
		return
	}

	count := this.SlotRepository.GetRemainingSlotByBuildingId(request.BuildingID)
	if count < 1 {
		helpers.HttpError(c, "Parking full", http.StatusBadRequest)
		return
	}

	parking := this.ParkingRepository.SearchParkingPlace(request.BuildingID)

	transaction := this.ParkingRepository.Create(request, parking)

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": &responses.TransactionStartResponse{
			PlatNo:   transaction.PlatNo,
			SlotName: parking.SlotName,
			Code:     transaction.Code,
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
// @Router /transactions/end/ [post]
func (this *TransactionsController) End(c *gin.Context) {
	var request requests.TransactionEndRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		helpers.HttpError(c, "Wrong Format", http.StatusBadRequest)
		return
	}

	transaction, err := this.ParkingRepository.GetByCode(request.Code)
	if err != nil {
		helpers.HttpError(c, "Transaction not found", http.StatusNotFound)
		return
	}

	transaction = this.ParkingRepository.EndParking(transaction)

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   &responses.TransactionEndResponse{Total: transaction.Total},
	})
}
