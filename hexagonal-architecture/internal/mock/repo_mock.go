package mock

import (
	"myapp/internal/repository"
)

var _ repository.PaymentRepository = (*MockPaymentRepository)(nil)

// MockPaymentRepository is a mock implementation of the PaymentRepository interface.
type MockPaymentRepository struct {
	SaveFunc     func(paymentID string) error
	FindByIDFunc func(id string) (string, error)
}

func (m *MockPaymentRepository) Save(paymentID string) error {
	return m.SaveFunc(paymentID)
}

func (m *MockPaymentRepository) FindByID(id string) (string, error) {
	return m.FindByIDFunc(id)
}
