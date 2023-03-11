package loms

import (
	"context"
	"log"
	"route256/loms/internal/converter"
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

	orderItems := converter.LomsApiToOrderItemsListModel(req.Items)
	orderId, err := i.lomsService.CreateOrder(ctx, req.User, orderItems)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &desc.CreateOrderResponse{
		OrderID: orderId,
	}, nil
}
