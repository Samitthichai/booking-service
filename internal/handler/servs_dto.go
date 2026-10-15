package handler

type ServicesRequest struct {
	ServiceName     string  `json:"service_name" binding:"required"`
	Description     string  `json:"description"`
	DurationMinutes int     `json:"duration_minutes" binding:"required,gt=0"`
	Price           float64 `json:"price" binding:"gte=0"`
}

type ServicesResponse struct {
	ID              int64   `json:"id"`
	ServiceName     string  `json:"service_name"`
	Description     string  `json:"description"`
	DurationMinutes int     `json:"duration_minutes"`
	Price           float64 `json:"price"`
	IsActive        bool    `json:"is_active"`
}

type ServiceUpdateRequest struct {
	ServiceName     *string  `json:"service_name"`
	Description     *string  `json:"description"`
	DurationMinutes *int     `json:"duration_minutes"`
	Price           *float64 `json:"price"`
	IsActive        *bool    `json:"is_active"`
}
