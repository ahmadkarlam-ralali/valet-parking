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

type ParkingRepository struct {
	Db      *gorm.DB
	Parking Parking
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