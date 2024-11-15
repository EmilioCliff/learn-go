package repository

type PaymentRepository interface {
	Save(paymentID string) error
	FindByID(id string) (string, error)
}
