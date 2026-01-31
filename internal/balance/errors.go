package balance

import "errors"

var (
	ErrNotEnoughFunds = errors.New("not enough funds")
	ErrInvalidOrder   = errors.New("invalid order number")
)
