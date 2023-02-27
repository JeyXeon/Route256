package loms

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"log"
	desc "route256/loms/pkg"
)

func (i *Implementation) CancelOrder(ctx context.Context, req *desc.CancelOrderRequest) (*empty.Empty, error) {
	if req.OrderID == 0 {
		return &empty.Empty{}, ErrEmptyOrder
	}

	log.Printf("cancelOrder: %+v", req)

	return &empty.Empty{}, nil
}
