package loms

import (
	"context"
	"log"
	desc "route256/loms/pkg/loms"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/pkg/errors"
)

var (
	ErrCancelOrderEmptyOrder = errors.New("empty order")
)

func (i *Implementation) CancelOrder(ctx context.Context, req *desc.CancelOrderRequest) (*empty.Empty, error) {
	if req.OrderID == 0 {
		return &empty.Empty{}, ErrCancelOrderEmptyOrder
	}

	log.Printf("cancelOrder: %+v", req)

	return &empty.Empty{}, nil
}
