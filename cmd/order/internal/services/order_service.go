package services

import (
	"context"
	"errors"

	paymentGrpc "github.com/tricong1998/go-ecom/cmd/order/internal/gateway/payment/grpc"
	productGrpc "github.com/tricong1998/go-ecom/cmd/order/internal/gateway/product/grpc"
	userGrpc "github.com/tricong1998/go-ecom/cmd/order/internal/gateway/user/grpc"
	"github.com/tricong1998/go-ecom/cmd/order/internal/models"
	"github.com/tricong1998/go-ecom/cmd/order/internal/repository"
	"github.com/tricong1998/go-ecom/cmd/payment/pkg/pb"
	"github.com/tricong1998/go-ecom/cmd/user/pkg/dto"
	"github.com/tricong1998/go-ecom/pkg/rabbitmq"
)

const (
	Pending = "pending"
	Success = "success"
	Failed  = "failed"
)

type OrderService struct {
	OrderRepo            repository.IOrderRepository
	UserGrpcGateway      userGrpc.IUserGateway
	PaymentGrpcGateway   paymentGrpc.IPaymentGateway
	ProductGrpcGateway   productGrpc.IProductGateway
	CreateOrderPublisher rabbitmq.IPublisher
}

type IOrderService interface {
	CreateOrder(input *models.Order) error
	ReadOrder(id uint) (*models.Order, error)
	ListOrders(
		perPage, page int32,
		userId uint,
	) ([]models.Order, int64, error)
	UpdateOrder(user *models.Order) error
	DeleteOrder(id uint) error
}

func NewOrderService(
	userRepo repository.IOrderRepository,
	userGateway userGrpc.IUserGateway,
	paymentGateway paymentGrpc.IPaymentGateway,
	createOrderPublisher rabbitmq.IPublisher,
	productGateway productGrpc.IProductGateway,
) *OrderService {
	return &OrderService{userRepo, userGateway, paymentGateway, productGateway, createOrderPublisher}
}

func (us *OrderService) CreateOrder(order *models.Order) error {
	user, err := us.UserGrpcGateway.Get(context.Background(), uint(order.UserId))
	if err != nil {
		return err
	}
	product, err := us.ProductGrpcGateway.Get(context.Background(), uint(order.ProductId))
	if err != nil {
		return err
	}
	order.Username = user.Username
	order.Amount = uint(product.GetProduct().GetPrice()) * uint(order.ProductCount)
	order.Status = Pending
	err = us.OrderRepo.CreateOrder(order)
	if err != nil {
		return err
	}

	err = us.PaymentOrder(order)
	if err != nil {
		return err
	}

	return nil
}

func (us *OrderService) PaymentOrder(order *models.Order) error {
	payment, err := us.PaymentGrpcGateway.
		Create(context.Background(), &pb.CreatePaymentRequest{
			OrderId: uint64(order.ID),
			Amount:  uint64(order.Amount),
			Method:  "cash",
			UserId:  uint64(order.UserId),
		})
	if err != nil {
		return err
	}

	if payment.GetPayment().Status == "failed" {
		order.Status = Failed
		err = us.OrderRepo.UpdateOrderStatus(order.ID, Failed)
		if err != nil {
			return err
		}
	}

	success, err := us.ProductGrpcGateway.UpdateProductQuantity(context.Background(), uint(order.ProductId), uint(order.ProductCount))
	if err != nil {
		return err
	}
	if !success {
		return errors.New("failed to update product quantity")
	}

	order.Status = Success
	err = us.OrderRepo.UpdateOrderStatus(order.ID, Success)
	if err != nil {
		return err
	}

	createUserPoint := dto.CreateUserPoint{
		OrderId: order.ID,
		UserId:  uint(order.UserId),
		Amount:  uint(order.Amount),
	}
	err = us.CreateOrderPublisher.PublishMessage(createUserPoint)
	if err != nil {
		return err
	}

	return nil
}

func (us *OrderService) ReadOrder(id uint) (*models.Order, error) {
	user, err := us.OrderRepo.ReadOrder(id)
	return user, err
}

func (us *OrderService) ListOrders(
	perPage, page int32,
	userId uint,
) ([]models.Order, int64, error) {
	return us.OrderRepo.ListOrders(perPage, page, userId)
}

func (us *OrderService) UpdateOrder(user *models.Order) error {
	err := us.OrderRepo.UpdateOrder(user)
	return err
}

func (us *OrderService) DeleteOrder(id uint) error {
	return us.OrderRepo.DeleteOrder(id)
}
