package services_test

import (
	"errors"
	"testing"

	"myapp/internal/mock"
	"myapp/internal/services"
	"myapp/pkg"

	"github.com/stretchr/testify/assert"
)

func TestSavePayment_Success(t *testing.T) {
	mockPaymentRepo := &mock.MockPaymentRepository{
		SaveFunc: func(paymentID string) error {
			// Simulating successful save
			return nil
		},
		FindByIDFunc: func(id string) (string, error) {
			return "someData", nil
		},
	}

	mockPublisher := &mock.MockMessagePublisher{
		PublishPaymentEventFunc: func(paymentID string) error {
			// Simulating successful event publishing
			return nil
		},
	}

	paymentService := services.NewPaymentService(mockPaymentRepo, mockPublisher)

	err := paymentService.ProcessPayment()
	assert.NoError(t, err, "Expected no error, but got: %v", err)
}

func TestSavePayment_FailureOnSave(t *testing.T) {
	mockPaymentRepo := &mock.MockPaymentRepository{
		SaveFunc: func(paymentID string) error {
			// Simulating an error on saving
			return errors.New("save failed")
		},
	}

	mockPublisher := &mock.MockMessagePublisher{
		PublishPaymentEventFunc: func(paymentID string) error {
			// This should not be called due to save failure
			return nil
		},
	}

	paymentService := services.NewPaymentService(mockPaymentRepo, mockPublisher)
	err := paymentService.ProcessPayment()
	assert.Error(t, err, "Expected error but got none")
	assert.EqualError(t, err, pkg.ErrFailedToSavePayment.Error(), "Expected error message 'save failed', but got: %v", err)
}

func TestSavePayment_FailureOnPublishEvent(t *testing.T) {
	mockPaymentRepo := &mock.MockPaymentRepository{
		SaveFunc: func(paymentID string) error {
			// Simulating successful save
			return nil
		},
		FindByIDFunc: func(id string) (string, error) {
			return "someData", nil
		},
	}

	mockPublisher := &mock.MockMessagePublisher{
		PublishPaymentEventFunc: func(paymentID string) error {
			// Simulating error on publishing event
			return errors.New("publish failed")
		},
	}

	paymentService := services.NewPaymentService(mockPaymentRepo, mockPublisher)

	err := paymentService.ProcessPayment()
	assert.Error(t, err, "Expected error but got none")
	assert.EqualError(
		t,
		err,
		pkg.ErrFailedToPublishEvent.Error(),
		"Expected error message 'payment saved but event publish failed: publish failed', but got: %v",
		err,
	)
}
