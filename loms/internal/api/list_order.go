package loms

import (
	"context"
	"route256/loms/internal/converter"
	desc "route256/loms/pkg/loms"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrListOrderEmptyOrder = status.Error(codes.InvalidArgument, "empty order")
)

func (i *Implementation) ListOrder(ctx context.Context, req *desc.ListOrderRequest) (*desc.ListOrderResponse, error) {
	if req.OrderID == 0 {
		return nil, ErrListOrderEmptyOrder
	}

	order, err := i.lomsService.ListOrder(ctx, req.OrderID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &desc.ListOrderResponse{
		Status: string(order.Status),
		User:   order.User,
		Items:  converter.ModelToOrderItemsListLomsApi(order.Items),
	}, nil
}
