package checkout

import (
	"context"
	"log"
	desc "route256/checkout/pkg/checkout"

	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	ErrDeleteFromCartEmptyUser = errors.New("empty user")
	ErrDeleteFromCartEmptySKU  = errors.New("empty sku")
)

func (i *Implementation) DeleteFromCart(ctx context.Context, req *desc.DeleteFromCartRequest) (*emptypb.Empty, error) {
	if req.User == 0 {
		return &emptypb.Empty{}, ErrDeleteFromCartEmptyUser
	}
	if req.Sku == 0 {
		return &emptypb.Empty{}, ErrDeleteFromCartEmptySKU
	}

	log.Printf("deleteFromCart: %+v", req)

	err := i.checkoutService.DeleteFromCart(ctx, req.User, req.Sku, req.Count)
	if err != nil {
		return &emptypb.Empty{}, err
	}

	return &emptypb.Empty{}, nil
}
