package pkg

import "errors"

var (
	ErrProcessingPayment    = errors.New("error processing payment")
	ErrFailedToSavePayment  = errors.New("failed to save payment")
	ErrFailedToPublishEvent = errors.New("failed to publish payment event")
)
