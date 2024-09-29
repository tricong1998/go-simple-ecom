package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/tricong1998/go-ecom/cmd/product/pkg/pb"
)

type MockProductGateway struct {
	mock.Mock
}

func (m *MockProductGateway) Get(ctx context.Context, productId uint) (*pb.ReadProductResponse, error) {
	args := m.Called(ctx, productId)
	return args.Get(0).(*pb.ReadProductResponse), args.Error(1)
}

func (m *MockProductGateway) UpdateProductQuantity(ctx context.Context, productId uint, quantity uint) (bool, error) {
	args := m.Called(ctx, productId, quantity)
	return args.Bool(0), args.Error(1)
}
