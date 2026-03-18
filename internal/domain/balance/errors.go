package balance

import "errors"

var (
	ErrInsufficientFunds          = errors.New("insufficient funds")
	ErrIncorrectOrderNumberFormat = errors.New("incorrect order number format")
)
