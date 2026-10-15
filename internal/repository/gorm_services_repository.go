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

func (r *GormServicesRepository) GetAllServices(includeInactive bool) ([]model.Services, error) {
	var services []model.Services

	q := r.db.Model(&model.Services{})
	if !includeInactive {
		q = q.Where("is_active = ?", true)
	}
	if err := q.Find(&services).Error; err != nil {
		return nil, err
	}
	return services, nil
}

func (r *GormServicesRepository) GetServiceByID(id int64) (*model.Services, error) {
	var s model.Services
	if err := r.db.First(&s, id).Error; err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *GormServicesRepository) UpdateService(services *model.Services) error {
	return r.db.Save(services).Error
}

func (r *GormServicesRepository) DeleteService(id int64) error {
	return r.db.Delete(&model.Services{}, id).Error
}
