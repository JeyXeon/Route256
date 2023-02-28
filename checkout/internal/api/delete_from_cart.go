package checkout

import (
	"context"
	"log"
	desc "route256/checkout/pkg/checkout"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (i *Implementation) DeleteFromCart(ctx context.Context, req *desc.DeleteFromCartRequest) (*emptypb.Empty, error) {
	if req.User == 0 {
		return &emptypb.Empty{}, ErrEmptyUser
	}
	if req.Sku == 0 {
		return &emptypb.Empty{}, ErrEmptySKU
	}

	log.Printf("deleteFromCart: %+v", req)

	err := i.checkoutService.DeleteFromCart(ctx, req.User, req.Sku, req.Count)
	if err != nil {
		return &emptypb.Empty{}, err
	}

	return &emptypb.Empty{}, nil
}
