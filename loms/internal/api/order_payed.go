package loms

import (
	"context"
	"log"
	desc "route256/loms/pkg/loms"

	"github.com/golang/protobuf/ptypes/empty"
)

func (i *Implementation) OrderPayed(ctx context.Context, req *desc.OrderPayedRequest) (*empty.Empty, error) {
	if req.OrderID == 0 {
		return &empty.Empty{}, ErrEmptyOrder
	}

	log.Printf("orderPayed: %+v", req)

	return &empty.Empty{}, nil
}
