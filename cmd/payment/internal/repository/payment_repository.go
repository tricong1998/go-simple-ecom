package repository

import (
	"github.com/tricong1998/go-ecom/cmd/payment/pkg/models"
	"gorm.io/gorm"
)

type PaymentRepository struct {
	db *gorm.DB
}

type IPaymentRepository interface {
	CreatePayment(input *models.Payment) error
	ReadPayment(id uint) (*models.Payment, error)
	ListPayments(
		perPage, page int32,
		userId *uint,
	) ([]models.Payment, int64, error)
	UpdatePayment(input *models.Payment) error
	DeletePayment(id uint) error
}

func NewPaymentRepository(db *gorm.DB) *PaymentRepository {
	return &PaymentRepository{db}
}

func (paymentRepo *PaymentRepository) CreatePayment(input *models.Payment) error {
	return paymentRepo.db.Create(input).Error
}

func (paymentRepo *PaymentRepository) ReadPayment(id uint) (*models.Payment, error) {
	var payment *models.Payment
	err := paymentRepo.db.First(&payment, id).Error
	if err != nil {
		return nil, err
	}
	return payment, nil
}

func (paymentRepo *PaymentRepository) ListPayments(
	perPage, page int32,
	userId *uint,
) ([]models.Payment, int64, error) {
	var payments []models.Payment
	var total int64

	var query models.Payment
	if userId != nil {
		query.UserID = *userId
	}

	err := paymentRepo.db.Model(&models.Payment{}).Where(query).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	paymentRepo.db.Where(query).Find(&payments)

	return payments, total, nil
}

func (paymentRepo *PaymentRepository) UpdatePayment(input *models.Payment) error {
	return paymentRepo.db.Save(input).Error
}

func (paymentRepo *PaymentRepository) DeletePayment(id uint) error {
	return paymentRepo.db.Delete(&models.Payment{}, id).Error
}
