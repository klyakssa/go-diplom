package balance

import "errors"

var (
	ErrInsufficientFunds          = errors.New("insufficient funds")            // return when insufficient funds
	ErrIncorrectOrderNumberFormat = errors.New("incorrect order number format") // return when incorrect order number format
)
