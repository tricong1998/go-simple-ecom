package services

import (
	"context"

	"github.com/tricong1998/go-ecom/cmd/order/internal/gateway/user/grpc"
	"github.com/tricong1998/go-ecom/cmd/order/internal/models"
	"github.com/tricong1998/go-ecom/cmd/order/internal/repository"
	"github.com/tricong1998/go-ecom/cmd/user/pkg/dto"
	"github.com/tricong1998/go-ecom/pkg/rabbitmq"
)

type OrderService struct {
	OrderRepo            repository.IOrderRepository
	UserGrpcGateway      grpc.IUserGateway
	CreateOrderPublisher rabbitmq.IPublisher
}

type IOrderService interface {
	CreateOrder(input *models.Order) error
	ReadOrder(id uint) (*models.Order, error)
	ListOrders(
		perPage, page int32,
		username int,
	) ([]models.Order, int64, error)
	UpdateOrder(user *models.Order) error
	DeleteOrder(id uint) error
}

func NewOrderService(userRepo repository.IOrderRepository, userGateway grpc.IUserGateway, createOrderPublisher rabbitmq.IPublisher) *OrderService {
	return &OrderService{userRepo, userGateway, createOrderPublisher}
}

func (us *OrderService) CreateOrder(order *models.Order) error {
	user, err := us.UserGrpcGateway.Get(context.Background(), uint(order.UserId))
	if err != nil {
		return err
	}
	order.Username = user.Username
	// TODO: calculate amount
	order.Amount = 1

	err = us.OrderRepo.CreateOrder(order)
	if err != nil {
		return err
	}

	createUserPoint := dto.CreateUserPoint{
		OrderId: order.ID,
		UserId:  uint(user.Id),
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
	username int,
) ([]models.Order, int64, error) {
	return us.OrderRepo.ListOrders(perPage, page, username)
}

func (us *OrderService) UpdateOrder(user *models.Order) error {
	err := us.OrderRepo.UpdateOrder(user)
	return err
}

func (us *OrderService) DeleteOrder(id uint) error {
	return us.OrderRepo.DeleteOrder(id)
}
