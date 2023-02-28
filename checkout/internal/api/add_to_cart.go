package checkout

import (
	"context"
	"log"
	desc "route256/checkout/pkg/checkout"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (i *Implementation) AddToCart(ctx context.Context, req *desc.AddToCartRequest) (*emptypb.Empty, error) {
	if req.User == 0 {
		return &emptypb.Empty{}, ErrEmptyUser
	}
	if req.Sku == 0 {
		return &emptypb.Empty{}, ErrEmptySKU
	}

	log.Printf("addToCart: %+v", req)

	err := i.checkoutService.AddToCart(ctx, req.User, req.Sku, req.Count)
	if err != nil {
		return &emptypb.Empty{}, err
	}

	return &emptypb.Empty{}, nil
}
