package services

import (
	"myapp/internal/repository"
	"myapp/pkg"
)

type PaymentService interface {
	ProcessPayment() error
	GetPaymentByID(id string) (string, error)
}

type paymentService struct {
	repo      repository.PaymentRepository
	publisher repository.MessagePublisher
}

func NewPaymentService(repo repository.PaymentRepository, publisher repository.MessagePublisher) PaymentService {
	return &paymentService{
		repo:      repo,
		publisher: publisher,
	}
}

func (s *paymentService) ProcessPayment() error {
	// Business logic to process a payment
	paymentID := "12345" // Example payment ID

	if err := s.repo.Save(paymentID); err != nil {
		return pkg.ErrFailedToSavePayment
	}

	if err := s.publisher.PublishPaymentEvent(paymentID); err != nil {
		return pkg.ErrFailedToPublishEvent
	}

	return nil
}

func (s *paymentService) GetPaymentByID(id string) (string, error) {
	return s.repo.FindByID(id)
}
