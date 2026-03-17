package orders

import "errors"

var (
	ErrOrderAlreadyAddedByThisUser  = errors.New("order already added by this user")
	ErrOrderAlreadyAddedByOtherUser = errors.New("order already added by other user")
	ErrIncorrectOrderNumberFormat   = errors.New("incorrect order number format")
)
