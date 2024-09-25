package services

import (
	"github.com/tricong1998/go-ecom/cmd/payment/internal/repository"
	"github.com/tricong1998/go-ecom/cmd/payment/pkg/models"
	"golang.org/x/exp/rand"
)

const (
	PaymentStatusPending = "pending"
	PaymentStatusSuccess = "success"
	PaymentStatusFailed  = "failed"
)

type PaymentService struct {
	PaymentRepo repository.IPaymentRepository
}

type IPaymentService interface {
	CreatePayment(input *models.Payment) error
	ReadPayment(id uint) (*models.Payment, error)
	ListPayments(
		perPage, page int32,
		userId *uint,
	) ([]models.Payment, int64, error)
	UpdatePayment(payment *models.Payment) error
	DeletePayment(id uint) error
}

func NewPaymentService(paymentRepo repository.IPaymentRepository) *PaymentService {
	return &PaymentService{paymentRepo}
}

func (us *PaymentService) CreatePayment(payment *models.Payment) error {
	payment.Status = PaymentStatusPending
	err := us.PaymentRepo.CreatePayment(payment)
	if err != nil {
		return err
	}
	res, err := us.ExecutePayment(payment)
	if err != nil {
		payment.Status = PaymentStatusFailed
		payment.Error = err.Error()
	}
	if res {
		payment.Status = PaymentStatusSuccess
	} else {
		payment.Status = PaymentStatusFailed
	}
	err = us.PaymentRepo.UpdatePayment(payment)
	return err
}

// ExecutePayment is a mock function to simulate payment execution, failure rate is 10%
func (us *PaymentService) ExecutePayment(payment *models.Payment) (bool, error) {
	random := rand.Intn(9)
	if random > 0 {
		return true, nil
	}
	return false, nil
}

func (us *PaymentService) ReadPayment(id uint) (*models.Payment, error) {
	payment, err := us.PaymentRepo.ReadPayment(id)
	return payment, err
}

func (us *PaymentService) ListPayments(
	perPage, page int32,
	userID *uint,
) ([]models.Payment, int64, error) {
	return us.PaymentRepo.ListPayments(perPage, page, userID)
}

func (us *PaymentService) UpdatePayment(payment *models.Payment) error {
	err := us.PaymentRepo.UpdatePayment(payment)
	return err
}

func (us *PaymentService) DeletePayment(id uint) error {
	return us.PaymentRepo.DeletePayment(id)
}
