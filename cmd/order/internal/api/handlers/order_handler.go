package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tricong1998/go-ecom/cmd/order/internal/api/dto"
	"github.com/tricong1998/go-ecom/cmd/order/internal/models"
	"github.com/tricong1998/go-ecom/cmd/order/internal/services"
)

type OrderHandler struct {
	OrderService services.IOrderService
}

func NewOrderHandler(userService services.IOrderService) *OrderHandler {
	return &OrderHandler{userService}
}

func (userHandler *OrderHandler) CreateOrder(ctx *gin.Context) {
	var input dto.CreateOrderDto
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user := models.Order{
		ProductId: input.ProductId,
		UserId:    input.UserId,
	}
	if err := userHandler.OrderService.CreateOrder(&user); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, dto.ToOrderResponse(&user))
}

func (userHandler *OrderHandler) ReadOrder(ctx *gin.Context) {
	var readOrderRequest dto.ReadOrderRequest
	if err := ctx.ShouldBindUri(&readOrderRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := userHandler.OrderService.ReadOrder(uint(readOrderRequest.ID))
	if err != nil {
		err := fmt.Errorf("user not found: %d", readOrderRequest.ID)
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, dto.ToOrderResponse(user))
}

func (userHandler *OrderHandler) UpdateOrder(ctx *gin.Context) {
	var input dto.CreateOrderDto
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var readOrderRequest dto.ReadOrderRequest
	if err := ctx.ShouldBindUri(&readOrderRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user := models.Order{
		ProductId: input.ProductId,
		UserId:    input.UserId,
	}
	user.ID = readOrderRequest.ID
	if err := userHandler.OrderService.UpdateOrder(&user); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, dto.ToOrderResponse(&user))
}

func (userHandler *OrderHandler) ListOrders(ctx *gin.Context) {
	var req dto.ListOrderQuery
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	users, total, err := userHandler.OrderService.ListOrders(req.PerPage, req.Page, req.UserId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var usersResponse []dto.OrderResponse
	for _, v := range users {
		usersResponse = append(usersResponse, *dto.ToOrderResponse(&v))
	}

	ctx.JSON(http.StatusOK, dto.ListOrderResponse{
		Items: usersResponse,
		Metadata: dto.MetadataDto{
			Total:   total,
			Page:    req.Page,
			PerPage: req.PerPage,
		},
	})
}

func (userHandler *OrderHandler) DeleteOrder(ctx *gin.Context) {
	var readOrderRequest dto.ReadOrderRequest
	if err := ctx.ShouldBindUri(&readOrderRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := userHandler.OrderService.DeleteOrder(uint(readOrderRequest.ID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}
