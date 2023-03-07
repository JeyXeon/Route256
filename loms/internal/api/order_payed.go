package loms

import (
	"context"
	"log"
	desc "route256/loms/pkg/loms"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrOrderPayedEmptyOrder = status.Error(codes.InvalidArgument, "empty order")
)

func (i *Implementation) OrderPayed(ctx context.Context, req *desc.OrderPayedRequest) (*empty.Empty, error) {
	if req.OrderID == 0 {
		return &empty.Empty{}, ErrOrderPayedEmptyOrder
	}

	log.Printf("orderPayed: %+v", req)

	return &empty.Empty{}, nil
}
