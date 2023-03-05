package checkout

import (
	"context"
	"log"
	desc "route256/checkout/pkg/checkout"

	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	ErrAddToCartEmptyUser = errors.New("empty user")
	ErrAddToCartEmptySKU  = errors.New("empty sku")
)

func (i *Implementation) AddToCart(ctx context.Context, req *desc.AddToCartRequest) (*emptypb.Empty, error) {
	if req.User == 0 {
		return &emptypb.Empty{}, ErrAddToCartEmptyUser
	}
	if req.Sku == 0 {
		return &emptypb.Empty{}, ErrAddToCartEmptySKU
	}

	log.Printf("addToCart: %+v", req)

	err := i.checkoutService.AddToCart(ctx, req.User, req.Sku, req.Count)
	if err != nil {
		return &emptypb.Empty{}, err
	}

	return &emptypb.Empty{}, nil
}
