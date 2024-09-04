package repository

import (
	"github.com/tricong1998/go-ecom/internal/models"
	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

func (orderRepo *OrderRepository) createOrder(input *models.Order) error {
	return orderRepo.db.Create(input).Error
}
