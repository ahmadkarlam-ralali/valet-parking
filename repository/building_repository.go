package repository

import (
	"errors"
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

func (repository *BuildingRepository) Create(request requests.BuildingStoreRequest) (models.Building, error) {
	var count int
	repository.Db.Model(&models.Building{}).Where("name = ?", request.Name).Count(&count)
	if count > 0 {
		return models.Building{}, errors.New("duplicate building name")
	}

	building := models.Building{
		Name: request.Name,
	}

	repository.Db.Create(&building)
	return building, nil
}

func (repository *BuildingRepository) Update(buildingID uint, request requests.BuildingUpdateRequest) (models.Building, error) {
	var count int
	repository.Db.
		Model(&models.Building{}).
		Where("id = ?", buildingID).
		Count(&count)
	if count == 0 {
		return models.Building{}, errors.New("building not found")
	}

	repository.Db.
		Model(&models.Building{}).
		Where("name = ? and id <> ?", request.Name, buildingID).
		Count(&count)
	if count > 0 {
		return models.Building{}, errors.New("duplicate building name")
	}

	var building models.Building
	repository.Db.First(&building, "id = ?", buildingID)
	building.Name = request.Name
	repository.Db.Model(&building).Updates(building)
	return building, nil
}

func (repository *BuildingRepository) Delete(buildingID uint) error {
	var count int
	repository.Db.
		Model(&models.Building{}).
		Where("id = ?", buildingID).
		Count(&count)
	if count == 0 {
		return errors.New("building not found")
	}

	repository.Db.Delete(&models.Building{}, "id = ?", buildingID)
	return nil
}
