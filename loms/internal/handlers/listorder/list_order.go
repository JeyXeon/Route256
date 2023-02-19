package listorder

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

type Item struct {
	SKU   int32  `json:"sku"`
	Count uint16 `json:"count"`
}

type Response struct {
	Status string `json:"status"`
	User   int64  `json:"user"`
	Items  []Item `json:"items"`
}

type Handler struct{}

func New() *Handler {
	return &Handler{}
}

func (h *Handler) Handle(ctx context.Context, request Request) (Response, error) {
	log.Printf("listOrder: %+v", request)
	return Response{
		Status: "new",
		User:   5,
		Items: []Item{
			{
				SKU:   12,
				Count: 5,
			},
		},
	}, nil
}
