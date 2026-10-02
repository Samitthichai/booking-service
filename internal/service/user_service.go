package service

import (
	"booking-service/internal/model"
	"booking-service/internal/repository"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	CreateUser(email, password string) error
	LoginUser(email, password string) (string, error)
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

func (s *userService) LoginUser(email string, password string) (string, error) {
	dbUser, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbUser.PasswordHash), []byte(password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = dbUser.ID
	claims["email"] = dbUser.Email
	claims["role"] = dbUser.Role
	claims["exp"] = time.Now().Add(24 * time.Hour).Unix()

	signed, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", errors.New("invalid credentials")
	}
	return signed, nil
}
