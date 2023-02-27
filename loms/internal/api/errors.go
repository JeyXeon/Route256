package loms

import "github.com/pkg/errors"

var (
	ErrEmptyOrder = errors.New("empty order")
	ErrEmptyUser  = errors.New("empty user")
	ErrEmptyItems = errors.New("empty items")
	ErrEmptySKU   = errors.New("empty sku")
)
