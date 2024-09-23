package models

import "gorm.io/gorm"

type Order struct {
	*gorm.Model
	Id        int    `json:"id"`
	Status    string `json:"status"`
	ProductId int    `json:"product_id"`
	UserId    int    `json:"user_id"`
	Username  string `json:"username"`
}
