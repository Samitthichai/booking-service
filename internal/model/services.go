package model

import "time"

type Services struct {
	ID              int64     `gorm:"primaryKey"`
	Name            string    `gorm:"not null"`
	Description     string    `gorm:""`
	DurationMinutes int       `gorm:"not null"`
	Price           float64   `gorm:"not null"`
	IsActive        bool      `gorm:"not null;default:true"`
	CreatedAt       time.Time `gorm:"not null;default:now()"`
}
