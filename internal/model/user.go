package model

import "time"

const (
	RoleCustomer = "customer"
	RoleAdmin    = "admin"
)

type User struct {
	ID           int64     `gorm:"primaryKey" `
	Email        string    `gorm:"uniqueIndex;not null" `
	PasswordHash string    `gorm:"not null" `
	Role         string    `gorm:"not null;default:customer;check:role IN ('customer','admin')" `
	CreatedAt    time.Time `gorm:"not null;default:now()" `
}
