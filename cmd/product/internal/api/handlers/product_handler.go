package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tricong1998/go-ecom/cmd/product/internal/services"
	"github.com/tricong1998/go-ecom/cmd/product/pkg/dto"
	"github.com/tricong1998/go-ecom/cmd/product/pkg/models"
)

type ProductHandler struct {
	ProductService services.IProductService
}

func NewProductHandler(userService services.IProductService) *ProductHandler {
	return &ProductHandler{userService}
}

func (userHandler *ProductHandler) CreateProduct(ctx *gin.Context) {
	var input dto.CreateProductDto
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user := models.Product{
		Name:  input.Name,
		Price: input.Price,
	}
	if err := userHandler.ProductService.CreateProduct(&user); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, dto.ToProductResponse(&user))
}

func (userHandler *ProductHandler) ReadProduct(ctx *gin.Context) {
	var readProductRequest dto.ReadProductRequest
	if err := ctx.ShouldBindUri(&readProductRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := userHandler.ProductService.ReadProduct(uint(readProductRequest.ID))
	if err != nil {
		err := fmt.Errorf("user not found: %d", readProductRequest.ID)
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, dto.ToProductResponse(user))
}

func (userHandler *ProductHandler) UpdateProduct(ctx *gin.Context) {
	var input dto.CreateProductDto
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var readProductRequest dto.ReadProductRequest
	if err := ctx.ShouldBindUri(&readProductRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user := models.Product{
		Name:  input.Name,
		Price: input.Price,
	}
	user.ID = readProductRequest.ID
	if err := userHandler.ProductService.UpdateProduct(&user); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, dto.ToProductResponse(&user))
}

func (userHandler *ProductHandler) ListProducts(ctx *gin.Context) {
	var req dto.ListProductQuery
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	users, total, err := userHandler.ProductService.ListProducts(req.PerPage, req.Page, req.Name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var usersResponse []dto.ProductResponse
	for _, v := range users {
		usersResponse = append(usersResponse, *dto.ToProductResponse(&v))
	}

	ctx.JSON(http.StatusOK, dto.ListProductResponse{
		Items: usersResponse,
		Metadata: dto.MetadataDto{
			Total:   total,
			Page:    req.Page,
			PerPage: req.PerPage,
		},
	})
}

func (userHandler *ProductHandler) DeleteProduct(ctx *gin.Context) {
	var readProductRequest dto.ReadProductRequest
	if err := ctx.ShouldBindUri(&readProductRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := userHandler.ProductService.DeleteProduct(uint(readProductRequest.ID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}
