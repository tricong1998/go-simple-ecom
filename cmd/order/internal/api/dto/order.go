package dto

import (
	"time"

	"github.com/tricong1998/go-ecom/cmd/order/internal/models"
)

type CreateOrderDto struct {
	ProductId    uint `json:"product_id" binding:"required"`
	UserId       uint `json:"user_id" binding:"required"`
	ProductCount uint `json:"product_count" binding:"required"`
}

type ReadOrderRequest struct {
	ID uint `uri:"id" binding:"required,min=1"`
}

type OrderResponse struct {
	ID           uint      `json:"id"`
	ProductId    uint      `json:"product_id"`
	UserId       uint      `json:"user_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	ProductCount uint      `json:"product_count"`
	Amount       uint      `json:"amount"`
}

type ListOrderQuery struct {
	UserId  uint  `form:"user_id"`
	Page    int32 `form:"page" binding:"required,min=1"`
	PerPage int32 `form:"per_page" binding:"required,min=5,max=10"`
}

type ListOrderResponse struct {
	Items    []OrderResponse `json:"items"`
	Metadata MetadataDto     `json:"metadata"`
}

func ToOrderResponse(user *models.Order) *OrderResponse {
	return &OrderResponse{
		ID:           user.ID,
		ProductId:    user.ProductId,
		UserId:       user.UserId,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		ProductCount: user.ProductCount,
		Amount:       user.Amount,
	}
}
