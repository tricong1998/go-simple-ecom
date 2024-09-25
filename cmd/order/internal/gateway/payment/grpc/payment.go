package grpc

import (
	"context"
	"fmt"
	"log"

	"github.com/tricong1998/go-ecom/cmd/payment/pkg/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type IPaymentGateway interface {
	Get(ctx context.Context, paymentId uint) (*pb.Payment, error)
	Create(ctx context.Context, payment *pb.CreatePaymentRequest) (*pb.CreatePaymentResponse, error)
}

type PaymentGateway struct {
	host string
	port string
}

func New(host string, port string) *PaymentGateway {
	return &PaymentGateway{host, port}
}

func (g *PaymentGateway) Get(ctx context.Context, paymentId uint) (*pb.Payment, error) {
	address := fmt.Sprintf("%s:%s", g.host, g.port)

	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := pb.NewPaymentGrpcClient(conn)
	resp, err := client.ReadPayment(ctx, &pb.ReadPaymentRequest{Id: (uint64)(paymentId)})
	if err != nil {
		log.Println("Error getting payment:", err)
		return nil, err
	}
	return resp, nil
}

func (g *PaymentGateway) Create(ctx context.Context, payment *pb.CreatePaymentRequest) (*pb.CreatePaymentResponse, error) {
	address := fmt.Sprintf("%s:%s", g.host, g.port)

	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := pb.NewPaymentGrpcClient(conn)
	resp, err := client.CreatePayment(ctx, &pb.CreatePaymentRequest{
		OrderId: payment.OrderId,
		UserId:  payment.UserId,
		Amount:  payment.Amount,
		Method:  payment.Method,
	})
	if err != nil {
		log.Println("Error creating payment:", err)
		return nil, err
	}
	return resp, nil
}
