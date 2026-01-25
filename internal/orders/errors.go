package orders

import "errors"

var (
	ErrEmptyBody     = errors.New("empty request body")
	ErrInvalidNumber = errors.New("invalid order number")
	ErrAlreadyMine   = errors.New("order already uploaded by this user")
	ErrAlreadyExists = errors.New("order already uploaded by another user")
)
