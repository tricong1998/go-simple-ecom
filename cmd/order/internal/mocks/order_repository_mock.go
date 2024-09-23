package mocks

import (
	"github.com/stretchr/testify/mock"
	"github.com/tricong1998/go-ecom/cmd/order/internal/models"
)

type MockOrderRepository struct {
	mock.Mock
}

func (m *MockOrderRepository) CreateOrder(user *models.Order) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockOrderRepository) ReadOrder(id uint) (*models.Order, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Order), args.Error(1)
}

func (m *MockOrderRepository) ListOrders(
	perPage, page int32,
	username int,
) ([]models.Order, int64, error) {
	args := m.Called(perPage, page, username)
	return args.Get(0).([]models.Order), args.Get(1).(int64), args.Error(2)
}

func (m *MockOrderRepository) UpdateOrder(user *models.Order) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockOrderRepository) DeleteOrder(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}
