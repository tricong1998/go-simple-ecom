package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/tricong1998/go-ecom/cmd/payment/pkg/pb"
)

type MockPaymentGateway struct {
	mock.Mock
}

func (m *MockPaymentGateway) Get(ctx context.Context, userId uint) (*pb.Payment, error) {
	args := m.Called(ctx, userId)
	return args.Get(0).(*pb.Payment), args.Error(1)
}

func (m *MockPaymentGateway) Create(ctx context.Context, payment *pb.CreatePaymentRequest) (*pb.CreatePaymentResponse, error) {
	args := m.Called(ctx, payment)
	return args.Get(0).(*pb.CreatePaymentResponse), args.Error(1)
}
