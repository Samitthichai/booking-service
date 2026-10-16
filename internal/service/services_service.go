package service

import (
	"booking-service/internal/model"
	"booking-service/internal/repository"
	"errors"
)

type ServsService interface {
	CreateService(name, description string, price float64, durationTime int) error
	GetAllServices(includeInactive bool) ([]model.Services, error)
	GetService(id int64) (*model.Services, error)
	UpdateService(id int64, patch ServiceUpdatePatch) (*model.Services, error)
	DeleteService(id int64) error
}

type ServiceUpdatePatch struct {
	Name            *string
	Description     *string
	DurationMinutes *int
	Price           *float64
	IsActive        *bool
}

type servsService struct {
	servRepo repository.ServicesRepository
}

func NewServsService(servRepo repository.ServicesRepository) ServsService {
	return &servsService{servRepo: servRepo}
}

var ErrInvalidService = errors.New("invalid service name, duration, or price")
var ErrNotFoundServices = errors.New("Not found Services")
var ErrServiceNotFound = errors.New("service not found")

func (s *servsService) CreateService(name, description string, price float64, durationTime int) error {
	if name == "" || durationTime <= 0 || price < 0 {
		return ErrInvalidService
	}
	servs := &model.Services{
		Name:            name,
		Description:     description,
		Price:           price,
		DurationMinutes: durationTime,
		IsActive:        true,
	}

	return s.servRepo.CreateService(servs)
}

func (s *servsService) GetAllServices(includeInactive bool) ([]model.Services, error) {
	return s.servRepo.GetAllServices(includeInactive)
}

func (s *servsService) GetService(id int64) (*model.Services, error) {
	svc, err := s.servRepo.GetServiceByID(id)
	if err != nil {
		return nil, ErrServiceNotFound
	}
	return svc, nil
}

func (s *servsService) UpdateService(id int64, patch ServiceUpdatePatch) (*model.Services, error) {
	svc, err := s.servRepo.GetServiceByID(id)
	if err != nil {
		return nil, ErrServiceNotFound
	}

	if patch.Name != nil {
		svc.Name = *patch.Name
	}
	if patch.Description != nil {
		svc.Description = *patch.Description
	}
	if patch.DurationMinutes != nil {
		svc.DurationMinutes = *patch.DurationMinutes
	}
	if patch.Price != nil {
		svc.Price = *patch.Price
	}
	if patch.IsActive != nil {
		svc.IsActive = *patch.IsActive
	}

	if svc.Name == "" || svc.DurationMinutes <= 0 || svc.Price < 0 {
		return nil, ErrInvalidService
	}

	if err := s.servRepo.UpdateService(svc); err != nil {
		return nil, err
	}
	return svc, nil
}

func (s *servsService) DeleteService(id int64) error {
	if _, err := s.servRepo.GetServiceByID(id); err != nil {
		return ErrServiceNotFound
	}
	return s.servRepo.DeleteService(id)
}
