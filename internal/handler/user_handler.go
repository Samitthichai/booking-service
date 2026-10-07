package handler

import (
	"errors"

	"booking-service/internal/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) RegisterUser(c *gin.Context) {
	var req AuthenticationRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"message": "Invalid request", "status": 400})
		return
	}
	email := req.Email
	password := req.Password

	if err := h.userService.CreateUser(email, password); err != nil {
		if errors.Is(err, service.ErrEmailExists) {
			c.JSON(409, gin.H{"error": err.Error()})
			return
		}
		c.JSON(500, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(201, gin.H{
		"message": "Create Success",
		"status":  201,
	})

}

func (h *UserHandler) LoginUser(c *gin.Context) {
	var req AuthenticationRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"message": "Invalid request", "status": 400})
		return
	}

	token, err := h.userService.LoginUser(req.Email, req.Password)
	if err != nil {
		c.JSON(401, gin.H{"message": "Invalid credentials", "status": 401})
		return
	}

	c.SetCookie("token", token, 3600*24, "/", "", true, true)

	c.JSON(200, gin.H{
		"message": "Login Success",
		"status":  200,
	})
}
