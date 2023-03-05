package loms

import (
	"context"
	"log"
	desc "route256/loms/pkg/loms"

	"github.com/pkg/errors"
)

var (
	ErrListOrderEmptyOrder = errors.New("empty order")
)

func (i *Implementation) ListOrder(ctx context.Context, req *desc.ListOrderRequest) (*desc.ListOrderResponse, error) {
	if req.OrderID == 0 {
		return nil, ErrListOrderEmptyOrder
	}

	log.Printf("listOrder: %+v", req)

	return &desc.ListOrderResponse{
		Status: "new",
		User:   5,
		Items: []*desc.OrderItem{
			{
				Sku:   12,
				Count: 5,
			},
		},
	}, nil
}
