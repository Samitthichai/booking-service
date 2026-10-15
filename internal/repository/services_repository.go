package repository

import "booking-service/internal/model"

type ServicesRepository interface {
	CreateService(services *model.Services) error
	GetAllServices(includeInactive bool) ([]model.Services, error)
	GetServiceByID(id int64) (*model.Services, error)
	UpdateService(services *model.Services) error
	DeleteService(id int64) error
}
