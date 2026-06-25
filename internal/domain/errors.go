package domain

import "errors"

var (
	ErrSubscriptionNotFound = errors.New("subscription not found")
	ErrInvalidInput        = errors.New("invalid input")
)