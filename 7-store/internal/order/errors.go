package order

import "errors"

var (
	ErrOrderNotFound = errors.New("order not found")
	ErrNotMatchOrderId = errors.New("order id not match")
)