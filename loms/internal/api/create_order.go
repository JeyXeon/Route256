package loms

import (
	"context"
	"log"
	desc "route256/loms/pkg/loms"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrCreateOrderEmptyUser  = status.Error(codes.InvalidArgument, "empty user")
	ErrCreateOrderEmptyItems = status.Error(codes.InvalidArgument, "empty items")
)

func (i *Implementation) CreateOrder(ctx context.Context, req *desc.CreateOrderRequest) (*desc.CreateOrderResponse, error) {
	if req.User == 0 {
		return nil, ErrCreateOrderEmptyUser
	}
	if len(req.Items) == 0 {
		return nil, ErrCreateOrderEmptyItems
	}

	log.Printf("createOrder: %+v", req)

	return &desc.CreateOrderResponse{
		OrderID: 5,
	}, nil
}
