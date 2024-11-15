package handlers

import (
	"myapp/internal/services"
	"myapp/pkg"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	service services.PaymentService
}

func NewPaymentHandler(service services.PaymentService) *PaymentHandler {
	return &PaymentHandler{service: service}
}

func (h *PaymentHandler) CreatePayment(c *gin.Context) {
	err := h.service.ProcessPayment()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": pkg.ErrProcessingPayment.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Payment processed successfully"})
}

func (h *PaymentHandler) GetPayment(c *gin.Context) {
	paymentID := c.Param("id")

	payment, err := h.service.GetPaymentByID(paymentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, payment)
}
