package error

import "errors"

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserExists        = errors.New("user already exists")
	ErrValidation        = errors.New("validation error")
	ErrCredential        = errors.New("invalid credential")
	ErrMethodNotAllowed  = errors.New("method not allowed")
	ErrUnauthorized      = errors.New("unauthorized")
	ErrInternalServer    = errors.New("internal server error")
	ErrNotFound          = errors.New("not found")
	ErrRouteNotFound     = errors.New("route not found")
	ErrInsufficientStock = errors.New("insufficient stock")
)
