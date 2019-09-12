package repository

import (
	"github.com/ahmadkarlam-ralali/valet-parking/models"
	"github.com/ahmadkarlam-ralali/valet-parking/requests"
	"github.com/jinzhu/gorm"
)

type BuildingRepository struct {
	Db *gorm.DB
}

func (repository *BuildingRepository) GetAll() []models.Building {
	var buildings []models.Building
	repository.Db.Find(&buildings)
	return buildings
}

func (repository *BuildingRepository) Get(buildingID uint) (models.Building, error) {
	var building models.Building
	result := repository.Db.First(&building, "id = ?", buildingID)

	return building, result.Error
}

func (repository *BuildingRepository) Create(request requests.BuildingStoreRequest) {
	repository.Db.Create(&models.Building{
		Name: request.Name,
	})
}

func (repository *BuildingRepository) Update(buildingID uint, request requests.BuildingUpdateRequest) {
	var building models.Building
	repository.Db.First(&building, "id = ?", buildingID)
	building.Name = request.Name
	repository.Db.Model(&building).Updates(building)
}

func (repository *BuildingRepository) Delete(buildingID uint) {
	repository.Db.Delete(&models.Building{}, "id = ?", buildingID)
}
