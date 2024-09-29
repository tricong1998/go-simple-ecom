package dto

import (
	"time"

	"github.com/tricong1998/go-ecom/cmd/product/pkg/models"
)

type CreateProductDto struct {
	Name     string `json:"name" binding:"required"`
	Price    uint   `json:"price" binding:"required,min=1"`
	Quantity uint   `json:"quantity" binding:"required,min=1"`
}

type ReadProductRequest struct {
	ID uint `uri:"id" binding:"required,min=1"`
}

type ProductResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Price     uint      `json:"price"`
	Quantity  uint      `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ListProductQuery struct {
	Name    *string `form:"name"`
	Page    int32   `form:"page" binding:"required,min=1"`
	PerPage int32   `form:"per_page" binding:"required,min=5,max=10"`
}

type ListProductResponse struct {
	Items    []ProductResponse `json:"items"`
	Metadata MetadataDto       `json:"metadata"`
}

func ToProductResponse(user *models.Product) *ProductResponse {
	return &ProductResponse{
		ID:        user.ID,
		Name:      user.Name,
		Price:     user.Price,
		Quantity:  user.Quantity,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
