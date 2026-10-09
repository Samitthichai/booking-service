package model

import "time"

type Services struct {
	ID              int64     `gorm:"primaryKey"`
	Name            string    `gorm:"not null"`
	DurationMinutes int       `gorm:"not null"`
	Price           float64   `gorm:"not null"`
	CreatedAt       time.Time `gorm:"not null;default:now()"`
}
