package mock

import (
	"myapp/internal/repository"
)

var _ repository.MessagePublisher = (*MockMessagePublisher)(nil)

type MockMessagePublisher struct {
	PublishPaymentEventFunc func(paymentID string) error
}

func (m *MockMessagePublisher) PublishPaymentEvent(paymentID string) error {
	return m.PublishPaymentEventFunc(paymentID)
}
