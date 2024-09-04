package models

import "gorm.io/gorm"

type Order struct {
	*gorm.Model
	Id       int       `json:"id"`
	Status   string    `json:"status"`
	Username string    `json:"username"`
	User     User      `json:"user"`
	Products []Product `json:"products" gorm:"many2many:user_languages;"`
}
