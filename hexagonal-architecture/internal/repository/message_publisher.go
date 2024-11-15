package repository

type MessagePublisher interface {
	PublishPaymentEvent(paymentID string) error
}
