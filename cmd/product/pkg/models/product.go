package models

import "gorm.io/gorm"

type Product struct {
	*gorm.Model
	Name   string  `json:"name"`
	Amount int     `json:"amount"`
	Price  float32 `json:"price"`
}
