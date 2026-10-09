package service

import (
	"booking-service/internal/model"
	"booking-service/internal/repository"
	"errors"
)

type ServsService interface {
	CreateService(name string, price float64, durationTime int) error
	GetAllServices() ([]model.Services, error)
}
type servsService struct {
	servRepo repository.ServicesRepository
}

func NewServsService(servRepo repository.ServicesRepository) ServsService {
	return &servsService{servRepo: servRepo}
}

var ErrInvalidService = errors.New("invalid service name, duration, or price")
var ErrNotFoundServices = errors.New("Not found Services")

func (s *servsService) CreateService(name string, price float64, durationTime int) error {
	if name == "" || durationTime <= 0 || price < 0 {
		return ErrInvalidService
	}
	servs := &model.Services{
		Name:            name,
		Price:           price,
		DurationMinutes: durationTime,
	}

	return s.servRepo.CreateService(servs)
}

func (s *servsService) GetAllServices() ([]model.Services, error) {
	return s.servRepo.GetAllServices()
}
