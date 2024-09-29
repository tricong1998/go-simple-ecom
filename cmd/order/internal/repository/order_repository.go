package repository

import (
	"github.com/tricong1998/go-ecom/cmd/order/internal/models"
	"gorm.io/gorm"
)

type OrderRepository struct {
	DB *gorm.DB
}

type IOrderRepository interface {
	CreateOrder(input *models.Order) error
	ReadOrder(id uint) (*models.Order, error)
	ListOrders(
		perPage, page int32,
		userId uint,
	) ([]models.Order, int64, error)
	UpdateOrder(input *models.Order) error
	DeleteOrder(id uint) error
	UpdateOrderStatus(orderId uint, status string) error
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db}
}

func (orderRepo *OrderRepository) CreateOrder(input *models.Order) error {
	return orderRepo.DB.Create(input).Error
}

func (userRepo *OrderRepository) ReadOrder(id uint) (*models.Order, error) {
	var user *models.Order
	err := userRepo.DB.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (userRepo *OrderRepository) ListOrders(
	perPage, page int32,
	userId uint,
) ([]models.Order, int64, error) {
	var users []models.Order
	var total int64

	var query models.Order
	if userId != 0 {
		query.UserId = userId
	}

	err := userRepo.DB.Model(&models.Order{}).Where(query).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	userRepo.DB.Where(query).Find(&users)

	return users, total, nil
}

func (userRepo *OrderRepository) UpdateOrder(input *models.Order) error {
	return userRepo.DB.Save(input).Error
}

func (userRepo *OrderRepository) DeleteOrder(id uint) error {
	return userRepo.DB.Delete(&models.Order{}, id).Error
}

func (userRepo *OrderRepository) UpdateOrderStatus(orderId uint, status string) error {
	return userRepo.DB.Model(&models.Order{}).Where("id = ?", orderId).Update("status", status).Error
}
