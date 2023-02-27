package createorder

import (
	"context"
	"errors"
	"log"
)

type Item struct {
	SKU   int32  `json:"sku"`
	Count uint16 `json:"count"`
}

type Request struct {
	User  int64  `json:"user"`
	Items []Item `json:"items"`
}

var (
	ErrEmptyUser  = errors.New("empty user")
	ErrEmptyItems = errors.New("empty items")
)

func (r Request) Validate() error {
	if r.User == 0 {
		return ErrEmptyUser
	}
	if len(r.Items) == 0 {
		return ErrEmptyItems
	}
	return nil
}

type Response struct {
	OrderId int64 `json:"orderId"`
}

type Handler struct{}

func New() *Handler {
	return &Handler{}
}

func (h *Handler) Handle(ctx context.Context, request Request) (Response, error) {
	log.Printf("createOrder: %+v", request)
	return Response{
		OrderId: 5,
	}, nil
}
