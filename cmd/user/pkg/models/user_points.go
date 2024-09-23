package models

import (
	"time"

	"gorm.io/gorm"
)

type UserPoint struct {
	gorm.Model
	OrderId    uint      `json:"order_id" gorm:"unique"`
	UserId     uint      `json:"user_id" `
	Point      uint      `json:"point"`
	ExpiryTime time.Time `json:"expiry_time"`
}
