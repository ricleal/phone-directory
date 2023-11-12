package entities

import "errors"

var (
	ErrInvalidPhoneNumber = errors.New("invalid phone number")
	ErrInvalidAddress     = errors.New("invalid address")
	ErrInvalidName        = errors.New("invalid name")
)
