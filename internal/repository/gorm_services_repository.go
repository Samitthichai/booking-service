package repository

import (
	"booking-service/internal/model"

	"gorm.io/gorm"
)

type GormServicesRepository struct {
	db *gorm.DB
}

func NewGormServicesRepository(db *gorm.DB) *GormServicesRepository {
	return &GormServicesRepository{db: db}
}

func (r *GormServicesRepository) CreateService(services *model.Services) error {
	return r.db.Create(services).Error

}
