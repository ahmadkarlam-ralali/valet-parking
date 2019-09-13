package v1

import (
	"github.com/ahmadkarlam-ralali/valet-parking/repository"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ReportsController struct {
	ParkingRepository repository.ParkingRepository
}

// Report total parking godoc
// @Summary Report total parking
// @Description report total parking
// @Tags Reports
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param type query string true "type (year, month)"
// @Param date query string true "date (2019, 2019-09)"
// @Success 200 {string} string "Ok"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /reports/total-parking/ [get]
func (this *ReportsController) GetTotalParking(c *gin.Context) {
	reports := this.ParkingRepository.GetTotalParking(c.Query("type"), c.Query("date"))
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   reports,
	})
}

// Report total income godoc
// @Summary Report total income
// @Description report total income
// @Tags Reports
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param date query string true "date"
// @Success 200 {string} string "Ok"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /reports/total-income/ [get]
func (this *ReportsController) GetIncomeParking(c *gin.Context) {
	reports := this.ParkingRepository.GetTotalIncomeByMonth(c.Query("date"))
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   reports,
	})
}
