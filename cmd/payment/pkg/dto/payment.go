package dto

import (
	"time"

	"github.com/tricong1998/go-ecom/cmd/payment/pkg/models"
)

type CreatePaymentDto struct {
	OrderId uint   `json:"order_id" binding:"required"`
	UserId  uint   `json:"user_id" binding:"required"`
	Amount  uint   `json:"amount" binding:"required"`
	Method  string `json:"method" binding:"required"`
}

type ReadPaymentRequest struct {
	ID uint `uri:"id" binding:"required,min=1"`
}

type PaymentResponse struct {
	ID        uint      `json:"id"`
	OrderId   uint      `json:"order_id"`
	UserId    uint      `json:"user_id"`
	Amount    uint      `json:"amount"`
	Method    string    `json:"method"`
	Status    string    `json:"status"`
	Error     string    `json:"error"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ListPaymentQuery struct {
	UserId  *uint `form:"user_id"`
	Page    int32 `form:"page" binding:"required,min=1"`
	PerPage int32 `form:"per_page" binding:"required,min=5,max=10"`
}

type ListPaymentResponse struct {
	Items    []PaymentResponse `json:"items"`
	Metadata MetadataDto       `json:"metadata"`
}

func ToPaymentResponse(payment *models.Payment) *PaymentResponse {
	return &PaymentResponse{
		ID:        payment.ID,
		OrderId:   payment.OrderID,
		UserId:    payment.UserID,
		Amount:    payment.Amount,
		Method:    payment.Method,
		Status:    payment.Status,
		CreatedAt: payment.CreatedAt,
		UpdatedAt: payment.UpdatedAt,
		Error:     payment.Error,
	}
}
