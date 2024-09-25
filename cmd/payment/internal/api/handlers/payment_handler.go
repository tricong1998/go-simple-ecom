package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tricong1998/go-ecom/cmd/payment/internal/services"
	"github.com/tricong1998/go-ecom/cmd/payment/pkg/dto"
	"github.com/tricong1998/go-ecom/cmd/payment/pkg/models"
)

type PaymentHandler struct {
	PaymentService services.IPaymentService
}

func NewPaymentHandler(paymentService services.IPaymentService) *PaymentHandler {
	return &PaymentHandler{paymentService}
}

func (paymentHandler *PaymentHandler) CreatePayment(ctx *gin.Context) {
	var input dto.CreatePaymentDto
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payment := models.Payment{
		OrderID: input.OrderId,
		UserID:  input.UserId,
		Amount:  input.Amount,
		Method:  input.Method,
	}
	if err := paymentHandler.PaymentService.CreatePayment(&payment); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, dto.ToPaymentResponse(&payment))
}

func (paymentHandler *PaymentHandler) ReadPayment(ctx *gin.Context) {
	var readPaymentRequest dto.ReadPaymentRequest
	if err := ctx.ShouldBindUri(&readPaymentRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payment, err := paymentHandler.PaymentService.ReadPayment(uint(readPaymentRequest.ID))
	if err != nil {
		err := fmt.Errorf("payment not found: %d", readPaymentRequest.ID)
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, dto.ToPaymentResponse(payment))
}

func (paymentHandler *PaymentHandler) UpdatePayment(ctx *gin.Context) {
	var input dto.CreatePaymentDto
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var readPaymentRequest dto.ReadPaymentRequest
	if err := ctx.ShouldBindUri(&readPaymentRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payment := models.Payment{
		OrderID: input.OrderId,
		UserID:  input.UserId,
		Amount:  input.Amount,
		Method:  input.Method,
	}
	payment.ID = readPaymentRequest.ID
	if err := paymentHandler.PaymentService.UpdatePayment(&payment); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, dto.ToPaymentResponse(&payment))
}

func (paymentHandler *PaymentHandler) ListPayments(ctx *gin.Context) {
	var req dto.ListPaymentQuery
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payments, total, err := paymentHandler.PaymentService.ListPayments(req.PerPage, req.Page, req.UserId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var paymentsResponse []dto.PaymentResponse
	for _, v := range payments {
		paymentsResponse = append(paymentsResponse, *dto.ToPaymentResponse(&v))
	}

	ctx.JSON(http.StatusOK, dto.ListPaymentResponse{
		Items: paymentsResponse,
		Metadata: dto.MetadataDto{
			Total:   total,
			Page:    req.Page,
			PerPage: req.PerPage,
		},
	})
}

func (paymentHandler *PaymentHandler) DeletePayment(ctx *gin.Context) {
	var readPaymentRequest dto.ReadPaymentRequest
	if err := ctx.ShouldBindUri(&readPaymentRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := paymentHandler.PaymentService.DeletePayment(uint(readPaymentRequest.ID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}
