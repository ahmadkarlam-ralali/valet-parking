package repository

import (
	"github.com/ahmadkarlam-ralali/valet-parking/models"
	"github.com/ahmadkarlam-ralali/valet-parking/requests"
	"github.com/jinzhu/gorm"
)

type SlotRepository struct {
	Db *gorm.DB
}

func (repository *SlotRepository) All(buildingID uint) []models.Slot {
	var slots []models.Slot
	repository.Db.Find(&slots, "building_id = ?", buildingID)
	return slots
}

func (repository *SlotRepository) Get(slotID uint) (models.Slot, error) {
	var slot models.Slot
	result := repository.Db.First(&slot, "id = ?", slotID)
	return slot, result.Error
}

func (repository *SlotRepository) Create(buildingID uint, request requests.SlotStoreRequest) models.Slot {
	slot := models.Slot{
		Name:       request.Name,
		BuildingID: buildingID,
		Total:      request.Total,
	}
	repository.Db.Create(&slot)
	return slot
}

func (repository *SlotRepository) Update(slot models.Slot, request requests.SlotUpdateRequest) models.Slot {
	slot.Name = request.Name
	slot.Total = request.Total
	repository.Db.Model(&slot).Updates(slot)
	return slot
}

func (repository *SlotRepository) GetTotalSlotOccupied(slotID uint) int {
	var count int
	repository.Db.Model(&models.Transaction{}).
		Where("slot_id = ? and end_at = '0000-00-00 00:00:00'", slotID).
		Count(&count)
	return count
}

func (repository *SlotRepository) DeleteByBuildingId(slotID uint, buildingID uint) {
	repository.Db.Delete(&models.Slot{}, "id = ? and building_id = ?", slotID, buildingID)
}

func (repository *SlotRepository) GetRemainingSlotByBuildingId(buildingID uint) int {
	type TotalData struct{ Total int }
	var sumSlot TotalData
	var countTransaction TotalData
	repository.Db.Table("slots").
		Select("sum(total) as total").
		Where("building_id = ?", buildingID).
		Scan(&sumSlot)
	repository.Db.Table("slots").
		Select("count(plat_no) as total").
		Joins("left join transactions on transactions.slot_id = slots.id").
		Where("building_id = ? and end_at = '0000-00-00 00:00:00'", buildingID).
		Scan(&countTransaction)

	return sumSlot.Total - countTransaction.Total
}
