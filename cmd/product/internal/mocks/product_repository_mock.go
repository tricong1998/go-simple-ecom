package mocks

import (
	"github.com/stretchr/testify/mock"
	"github.com/tricong1998/go-ecom/cmd/product/pkg/models"
)

type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) CreateProduct(user *models.Product) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockProductRepository) ReadProduct(id uint) (*models.Product, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Product), args.Error(1)
}

func (m *MockProductRepository) ListProducts(
	perPage, page int32,
	username *string,
) ([]models.Product, int64, error) {
	args := m.Called(perPage, page, username)
	return args.Get(0).([]models.Product), args.Get(1).(int64), args.Error(2)
}

func (m *MockProductRepository) UpdateProduct(user *models.Product) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockProductRepository) DeleteProduct(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockProductRepository) UpdateProductQuantity(productId, quantity uint) (bool, error) {
	args := m.Called(productId, quantity)
	return args.Get(0).(bool), args.Error(1)
}
