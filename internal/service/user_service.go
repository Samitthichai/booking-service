package service

import (
	"booking-service/internal/model"
	"booking-service/internal/repository"
	"crypto/subtle"
	"errors"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var ErrEmailExists = errors.New("email already exists")
var ErrInvalidCredentials = errors.New("email or password is incorrect")
var ErrInvalidInviteCode = errors.New("invalid invite code")

type UserService interface {
	CreateUser(email, password string) error
	CreateAdmin(email, password, inviteCode string) error
	LoginUser(email, password string) (string, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) CreateUser(email, password string) error {
	email = strings.ToLower(email)

	if _, err := s.userRepo.GetUserByEmail(email); err == nil {
		return ErrEmailExists
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

func (s *userService) CreateAdmin(email, password, inviteCode string) error {
	expected := os.Getenv("ADMIN_INVITE_CODE")
	if expected == "" || subtle.ConstantTimeCompare([]byte(inviteCode), []byte(expected)) != 1 {
		return ErrInvalidInviteCode
	}

	email = strings.ToLower(email)

	if _, err := s.userRepo.GetUserByEmail(email); err == nil {
		return ErrEmailExists
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	var user = &model.User{
		Email:        email,
		PasswordHash: string(hashPassword),
		Role:         model.RoleAdmin,
	}
	return s.userRepo.CreateUser(user)
}

func (s *userService) LoginUser(email string, password string) (string, error) {
	dbUser, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		return "", ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbUser.PasswordHash), []byte(password))
	if err != nil {
		return "", ErrInvalidCredentials
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = dbUser.ID
	claims["email"] = dbUser.Email
	claims["role"] = dbUser.Role
	claims["exp"] = time.Now().Add(24 * time.Hour).Unix()

	signed, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", ErrInvalidCredentials
	}
	return signed, nil
}
