package handlers

import (
	"myapp/internal/services"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, paymentService services.PaymentService) {
	paymentHandler := NewPaymentHandler(paymentService)

	router.POST("/payments", paymentHandler.CreatePayment)
	router.GET("/payments/:id", paymentHandler.GetPayment)
}
