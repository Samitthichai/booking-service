package repository

import "booking-service/internal/model"

type ServicesRepository interface {
	CreateService(services *model.Services) error
}
