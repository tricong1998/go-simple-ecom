package grpc_handler

import (
	"context"

	"github.com/tricong1998/go-ecom/cmd/payment/internal/services"
	"github.com/tricong1998/go-ecom/cmd/payment/pkg/models"
	"github.com/tricong1998/go-ecom/cmd/payment/pkg/pb"
)

type Server struct {
	PaymentService services.IPaymentService
	pb.UnimplementedPaymentGrpcServer
}

func NewServer(PaymentService services.IPaymentService) *Server {
	server := Server{
		PaymentService: PaymentService,
	}
	return &server
}

func (server *Server) ReadPayment(_ context.Context, input *pb.ReadPaymentRequest) (*pb.Payment, error) {
	payment, err := server.PaymentService.ReadPayment((uint)(input.GetId()))
	if err != nil {
		return nil, err
	}
	return &pb.Payment{
		Id:        uint64(payment.ID),
		OrderId:   uint64(payment.OrderID),
		UserId:    uint64(payment.UserID),
		Amount:    uint64(payment.Amount),
		Status:    payment.Status,
		CreatedAt: payment.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: payment.UpdatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

func (server *Server) CreatePayment(ctx context.Context, input *pb.CreatePaymentRequest) (*pb.CreatePaymentResponse, error) {
	payment := models.Payment{
		OrderID: (uint)(input.GetOrderId()),
		UserID:  (uint)(input.GetUserId()),
		Amount:  (uint)(input.GetAmount()),
		Method:  input.GetMethod(),
	}
	err := server.PaymentService.CreatePayment(&payment)
	if err != nil {
		return nil, err
	}
	return &pb.CreatePaymentResponse{
		Payment: &pb.Payment{
			Id:        uint64(payment.ID),
			OrderId:   uint64(payment.OrderID),
			UserId:    uint64(payment.UserID),
			Amount:    uint64(payment.Amount),
			Status:    payment.Status,
			CreatedAt: payment.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: payment.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	}, nil
}
