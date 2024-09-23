package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/tricong1998/go-ecom/cmd/user/pkg/pb"
)

type MockUserGateway struct {
	mock.Mock
}

func (m *MockUserGateway) Get(ctx context.Context, userId uint) (*pb.User, error) {
	args := m.Called(ctx, userId)
	return args.Get(0).(*pb.User), args.Error(1)
}
