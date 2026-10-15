package handler

import (
	"booking-service/internal/service"
	"errors"

	"github.com/gin-gonic/gin"
)

type ServiceHandler struct {
	servsService service.ServsService
}

func NewServsHandler(servsService service.ServsService) *ServiceHandler {
	return &ServiceHandler{servsService: servsService}
}

func (h *ServiceHandler) CreateServices(c *gin.Context) {
	var request ServicesRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	name := request.ServiceName
	price := request.Price
	durationTime := request.DurationMinutes

	if err := h.servsService.CreateService(name, price, durationTime); err != nil {
		if errors.Is(err, service.ErrInvalidService) {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(500, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(201, gin.H{
		"message": "Create Success",
	})
}

func (h *ServiceHandler) GetAllServices(c *gin.Context) {
	services, err := h.servsService.GetAllServices()

	if err != nil {
		if errors.Is(err, service.ErrNotFoundServices) {
			c.JSON(404, gin.H{"error": err.Error()})
			return
		}

		c.JSON(500, gin.H{"error": "internal server error"})
		return
	}
	response := make([]ServicesResponse, 0, len(services))
	for _, s := range services {
		response = append(response, ServicesResponse{
			ServiceName:     s.Name,
			DurationMinutes: s.DurationMinutes,
			Price:           s.Price,
		})
	}
	c.JSON(200, gin.H{
		"message": "Get Service Success",
		"data":    response,
	})
}
