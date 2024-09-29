package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username   string      `json:"username" gorm:"unique"`
	Password   string      `json:"password"`
	Role       string      `json:"role"`
	FullName   string      `json:"full_name"`
	UserPoints []UserPoint `json:"user_points"`
}

type Role string

const (
	AdminRole Role = "admin"
	UserRole  Role = "user"
)
