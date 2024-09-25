package mocks

import (
	"github.com/stretchr/testify/mock"
	"github.com/tricong1998/go-ecom/cmd/payment/pkg/models"
)

type MockPaymentRepository struct {
	mock.Mock
}

func (m *MockPaymentRepository) CreatePayment(payment *models.Payment) error {
	args := m.Called(payment)
	return args.Error(0)
}

func (m *MockPaymentRepository) ReadPayment(id uint) (*models.Payment, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Payment), args.Error(1)
}

func (m *MockPaymentRepository) ListPayments(
	perPage, page int32,
	userId *uint,
) ([]models.Payment, int64, error) {
	args := m.Called(perPage, page, userId)
	return args.Get(0).([]models.Payment), args.Get(1).(int64), args.Error(2)
}

func (m *MockPaymentRepository) UpdatePayment(payment *models.Payment) error {
	args := m.Called(payment)
	return args.Error(0)
}

func (m *MockPaymentRepository) DeletePayment(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}
