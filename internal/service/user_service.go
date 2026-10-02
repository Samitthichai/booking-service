package service

import (
	"booking-service/internal/model"
	"booking-service/internal/repository"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	CreateUser(email, password string) error
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) CreateUser(email, password string) error {

	if _, err := s.userRepo.GetUserByEmail(email); err == nil {
		return errors.New("email already exists")

	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	var user = &model.User{
		Email:        email,
		PasswordHash: string(hashPassword),
	}
	return s.userRepo.CreateUser(user)
}
