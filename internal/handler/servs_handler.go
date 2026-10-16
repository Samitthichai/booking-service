package handler

import (
	"booking-service/internal/model"
	"booking-service/internal/service"
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ServiceHandler struct {
	servsService service.ServsService
}

func NewServsHandler(servsService service.ServsService) *ServiceHandler {
	return &ServiceHandler{servsService: servsService}
}

func toServiceResponse(s model.Services) ServicesResponse {
	return ServicesResponse{
		ID:              s.ID,
		ServiceName:     s.Name,
		Description:     s.Description,
		DurationMinutes: s.DurationMinutes,
		Price:           s.Price,
		IsActive:        s.IsActive,
	}
}

func (h *ServiceHandler) CreateServices(c *gin.Context) {
	var request ServicesRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := h.servsService.CreateService(request.ServiceName, request.Description, request.Price, request.DurationMinutes); err != nil {
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

func (h *ServiceHandler) getAllServices(c *gin.Context, includeInactive bool) {
	services, err := h.servsService.GetAllServices(includeInactive)
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
		response = append(response, toServiceResponse(s))
	}
	c.JSON(200, gin.H{
		"message": "Get Service Success",
		"data":    response,
	})
}

func (h *ServiceHandler) GetAllServices(c *gin.Context) {
	h.getAllServices(c, false)
}

func (h *ServiceHandler) GetAllServicesAdmin(c *gin.Context) {
	h.getAllServices(c, true)
}

func (h *ServiceHandler) PatchService(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid service id"})
		return
	}

	var req ServiceUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	patch := service.ServiceUpdatePatch{
		Name:            req.ServiceName,
		Description:     req.Description,
		DurationMinutes: req.DurationMinutes,
		Price:           req.Price,
		IsActive:        req.IsActive,
	}

	updated, err := h.servsService.UpdateService(id, patch)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrServiceNotFound):
			c.JSON(404, gin.H{"error": err.Error()})
		case errors.Is(err, service.ErrInvalidService):
			c.JSON(400, gin.H{"error": err.Error()})
		default:
			c.JSON(500, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(200, gin.H{
		"message": "Update Success",
		"data":    toServiceResponse(*updated),
	})
}

func (h *ServiceHandler) DeleteService(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid service id"})
		return
	}

	if err := h.servsService.DeleteService(id); err != nil {
		if errors.Is(err, service.ErrServiceNotFound) {
			c.JSON(404, gin.H{"error": err.Error()})
			return
		}
		c.JSON(500, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(200, gin.H{"message": "Delete Success"})
}
