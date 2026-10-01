package model

import "time"

const (
	RoleCustomer = "customer"
	RoleAdmin    = "admin"
)

type User struct {
	ID           int64     `gorm:"primaryKey" json:"id"`
	Email        string    `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash string    `gorm:"not null" json:"-"`
	Role         string    `gorm:"not null;default:customer;check:role IN ('customer','admin')" json:"role"`
	CreatedAt    time.Time `gorm:"not null;default:now()" json:"created_at"`
}
