package api

import (
	"github.com/gin-gonic/gin"
	"github.com/tricong1998/go-ecom/cmd/payment/internal/api/handlers"
	"github.com/tricong1998/go-ecom/cmd/payment/internal/repository"
	"github.com/tricong1998/go-ecom/cmd/payment/internal/services"
	"gorm.io/gorm"
)

func SetupRoutes(routes *gin.Engine, db *gorm.DB) {
	paymentRepo := repository.NewPaymentRepository(db)
	paymentService := services.NewPaymentService(paymentRepo)
	paymentHandler := handlers.NewPaymentHandler(paymentService)

	paymentGroup := routes.Group("payments")
	{
		paymentGroup.POST("", paymentHandler.CreatePayment)
		paymentGroup.GET("/:id", paymentHandler.ReadPayment)
		paymentGroup.GET("", paymentHandler.ListPayments)
		paymentGroup.PUT("/:id", paymentHandler.UpdatePayment)
		paymentGroup.DELETE("/:id", paymentHandler.DeletePayment)
	}
}
