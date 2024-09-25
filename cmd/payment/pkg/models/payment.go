package models

import (
	"gorm.io/gorm"
)

type Payment struct {
	gorm.Model
	OrderID uint   `json:"order_id"`
	UserID  uint   `json:"user_id"`
	Amount  uint   `json:"amount"`
	Method  string `json:"method"`
	Status  string `json:"status"`
	Error   string `json:"error"`
}
