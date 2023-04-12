package loms

import (
	"context"
	desc "route256/loms/pkg/loms"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrCancelOrderEmptyOrder = status.Error(codes.InvalidArgument, "empty order")
)

func (i *Implementation) CancelOrder(ctx context.Context, req *desc.CancelOrderRequest) (*empty.Empty, error) {
	if req.OrderID == 0 {
		return &empty.Empty{}, ErrCancelOrderEmptyOrder
	}

	err := i.lomsService.CancelOrder(ctx, req.OrderID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &empty.Empty{}, nil
}
