package checkout

import (
	"context"
	"log"
	"route256/checkout/internal/converters"
	desc "route256/checkout/pkg/checkout"

	"github.com/pkg/errors"
)

var (
	ErrListCartEmptyUser = errors.New("empty user")
)

func (i *Implementation) ListCart(ctx context.Context, req *desc.ListCartRequest) (*desc.ListCartResponse, error) {
	if req.User == 0 {
		return nil, ErrListCartEmptyUser
	}

	log.Printf("listCart: %+v", req)

	cart, err := i.checkoutService.ListCart(ctx, req.User)
	if err != nil {
		return nil, err
	}

	return &desc.ListCartResponse{
		Items:      converters.ToProductListDesc(cart.Items),
		TotalPrice: cart.TotalPrice,
	}, nil
}
