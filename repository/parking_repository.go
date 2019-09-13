package repository

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/ahmadkarlam-ralali/valet-parking/models"
	"github.com/ahmadkarlam-ralali/valet-parking/requests"
	"github.com/jinzhu/gorm"
	"math"
	"time"
)

type Parking struct {
	SlotID         uint
	SlotName       string
	TotalParking   int
	TotalAvailable int
}

type ReportParking struct {
	Date  string `json:"date"`
	Total int    `json:"total"`
}

type ParkingRepository struct {
	Db      *gorm.DB
	Parking Parking
}

func (repository *ParkingRepository) GetAll() []models.Transaction {
	var transactions []models.Transaction
	repository.Db.Preload("Slot.Building").Find(&transactions)
	return transactions
}

func (repository *ParkingRepository) SearchParkingPlace(buildingID uint) Parking {
	rows, _ := repository.Db.Raw("select "+
		"s.id as slot_id, s.name as slot_name, "+
		"(select count(*) from transactions t where end_at = '0000-00-00 00:00:00' and t.slot_id = s.id) as total_parking, "+
		"s.total as total_available from slots s where building_id = ?", buildingID).
		Rows()

	var parking Parking
	for rows.Next() {
		_ = rows.Scan(&parking.SlotID, &parking.SlotName, &parking.TotalParking, &parking.TotalAvailable)

		if parking.TotalParking < parking.TotalAvailable {
			break
		}
	}

	return parking
}

func (repository *ParkingRepository) Create(request requests.TransactionStartRequest, parking Parking) models.Transaction {
	h := md5.New()
	h.Write([]byte(fmt.Sprintf("%s%s%d", request.PlatNo, time.Now(), parking.SlotID)))
	code := hex.EncodeToString(h.Sum(nil))

	transaction := models.Transaction{
		Code:    code,
		SlotId:  parking.SlotID,
		PlatNo:  request.PlatNo,
		StartAt: time.Now(),
	}
	repository.Db.Create(&transaction)

	return transaction
}

func (repository *ParkingRepository) GetByCode(Code string) (models.Transaction, error) {
	var transaction models.Transaction
	result := repository.Db.First(&transaction, "code = ?", Code)
	return transaction, result.Error
}

func (repository *ParkingRepository) EndParking(transaction models.Transaction) models.Transaction {
	transaction.EndAt = time.Now()
	duration := transaction.EndAt.Sub(transaction.StartAt).Hours()
	transaction.Total = int(1500 * math.Ceil(duration))
	repository.Db.Model(&transaction).Updates(transaction)
	return transaction
}

func (repository *ParkingRepository) IsPlatNoStillParking(PlatNo string) bool {
	var count int
	repository.Db.
		Model(&models.Transaction{}).
		Where("plat_no = ? and end_at = '0000-00-00 00:00:00'", PlatNo).Count(&count)
	return count > 0
}

func (repository *ParkingRepository) GetTotalParking(Type string, Date string) []ReportParking {
	var reports []ReportParking
	switch Type {
	case "month":
		reports = repository.GetTotalParkingByMonth(Date)
	case "year":
		reports = repository.GetTotalParkingByYear(Date)
	}
	return reports
}

func (repository *ParkingRepository) GetTotalParkingByYear(Date string) []ReportParking {
	var report []ReportParking
	repository.Db.Table("transactions").
		Select("month(start_at) as 'date', count(start_at) as 'total'").
		Where("end_at <> '0000-00-00 00:00:00' and date_format(start_at, '%Y') = ?", Date).
		Group("month(start_at)").
		Scan(&report)
	return report
}

func (repository *ParkingRepository) GetTotalParkingByMonth(Date string) []ReportParking {
	var report []ReportParking
	repository.Db.Table("transactions").
		Select("date(start_at) as 'date', count(start_at) as 'total'").
		Where("end_at <> '0000-00-00 00:00:00' and date_format(start_at, '%Y-%m') = ?", Date).
		Group("date(start_at)").
		Scan(&report)
	return report
}

func (repository *ParkingRepository) GetTotalIncome(Type string, Date string) []ReportParking {
	var reports []ReportParking
	switch Type {
	case "month":
		reports = repository.GetTotalIncomeByMonth(Date)
	case "year":
		reports = repository.GetTotalIncomeByYear(Date)
	}
	return reports
}

func (repository *ParkingRepository) GetTotalIncomeByYear(Date string) []ReportParking {
	var report []ReportParking
	repository.Db.Table("transactions").
		Select("month(start_at) as 'date', sum(total) as 'total'").
		Where("end_at <> '0000-00-00 00:00:00' and date_format(start_at, '%Y') = ?", Date).
		Group("month(start_at)").
		Scan(&report)
	return report
}

func (repository *ParkingRepository) GetTotalIncomeByMonth(Date string) []ReportParking {
	var report []ReportParking
	repository.Db.Table("transactions").
		Select("date(start_at) as 'date', sum(total) as 'total'").
		Where("end_at <> '0000-00-00 00:00:00' and date_format(start_at, '%Y-%m') = ?", Date).
		Group("date(start_at)").
		Scan(&report)
	return report
}
