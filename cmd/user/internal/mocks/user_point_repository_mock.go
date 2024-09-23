package mocks

import (
	"github.com/stretchr/testify/mock"
	"github.com/tricong1998/go-ecom/cmd/user/pkg/models"
)

type MockUserPointRepository struct {
	mock.Mock
}

func (m *MockUserPointRepository) CreateUserPoint(input *models.UserPoint) error {
	args := m.Called(input)
	return args.Error(0)
}

func (m *MockUserPointRepository) ReadUserPoint(id uint) (*models.UserPoint, error) {
	args := m.Called(id)
	return args.Get(0).(*models.UserPoint), args.Error(1)
}

func (m *MockUserPointRepository) UpdateUserPoint(input *models.UserPoint) error {
	args := m.Called(input)
	return args.Error(0)
}

func (m *MockUserPointRepository) DeleteUserPoint(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserPointRepository) ListUserPoints(
	perPage, page int32,
	userId *uint,
) ([]models.UserPoint, int64, error) {
	args := m.Called(perPage, page, userId)
	return args.Get(0).([]models.UserPoint), args.Get(1).(int64), args.Error(2)
}
