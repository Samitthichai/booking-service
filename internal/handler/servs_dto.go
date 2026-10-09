package handler

type ServicesRequest struct {
	ServiceName     string  `json:"service_name" binding:"required"`
	DurationMinutes int     `json:"duration_minutes" binding:"required,gt=0"`
	Price           float64 `json:"price" binding:"gte=0"`
}
