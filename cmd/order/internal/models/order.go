package models

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	Status       string `json:"status"`
	ProductId    uint   `json:"product_id"`
	UserId       uint   `json:"user_id"`
	Username     string `json:"username"`
	ProductCount uint   `json:"product_count"`
	Amount       uint   `json:"amount"`
}
