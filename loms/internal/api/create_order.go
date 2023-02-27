package loms

import (
	"context"
	"log"
	desc "route256/loms/pkg"
)

func (i *Implementation) CreateOrder(ctx context.Context, req *desc.CreateOrderRequest) (*desc.CreateOrderResponse, error) {
	if req.User == 0 {
		return nil, ErrEmptyUser
	}
	if len(req.Items) == 0 {
		return nil, ErrEmptyItems
	}

	log.Printf("createOrder: %+v", req)

	return &desc.CreateOrderResponse{
		OrderID: 5,
	}, nil
}
