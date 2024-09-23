package mocks

import "github.com/stretchr/testify/mock"

type MockRabbitPublisher struct {
	mock.Mock
}

func (m *MockRabbitPublisher) PublishMessage(msg interface{}) error {
	args := m.Called(msg)
	return args.Error(0)
}
