package cancelorder

import (
	"context"
	"errors"
	"log"
)

type Request struct {
	OrderId int64 `json:"orderId"`
}

var (
	ErrEmptyOrder = errors.New("empty order")
)

func (r Request) Validate() error {
	if r.OrderId == 0 {
		return ErrEmptyOrder
	}
	return nil
}

type Response struct {
	Test string `json:"test"`
}

type Handler struct{}

func New() *Handler {
	return &Handler{}
}

func (h *Handler) Handle(ctx context.Context, request Request) (Response, error) {
	log.Printf("cancelOrder: %+v", request)
	return Response{}, nil
}
