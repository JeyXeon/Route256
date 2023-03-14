package checkout

import (
	"context"
	"log"
	"route256/checkout/internal/converters"
	desc "route256/checkout/pkg/checkout"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrListCartEmptyUser = status.Error(codes.InvalidArgument, "empty user")
)

func (i *Implementation) ListCart(ctx context.Context, req *desc.ListCartRequest) (*desc.ListCartResponse, error) {
	if req.User == 0 {
		return nil, ErrListCartEmptyUser
	}

	log.Printf("listCart: %+v", req)

	cart, err := i.checkoutService.ListCart(ctx, req.User)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &desc.ListCartResponse{
		Items:      converters.ModelToProductListDesc(cart.Items),
		TotalPrice: cart.TotalPrice,
	}, nil
}
