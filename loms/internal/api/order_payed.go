package loms

import (
	"context"
	"log"
	desc "route256/loms/pkg/loms"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/pkg/errors"
)

var (
	ErrOrderPayedEmptyOrder = errors.New("empty order")
)

func (i *Implementation) OrderPayed(ctx context.Context, req *desc.OrderPayedRequest) (*empty.Empty, error) {
	if req.OrderID == 0 {
		return &empty.Empty{}, ErrOrderPayedEmptyOrder
	}

	log.Printf("orderPayed: %+v", req)

	return &empty.Empty{}, nil
}
