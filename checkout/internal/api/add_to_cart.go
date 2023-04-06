package checkout

import (
	"context"
	desc "route256/checkout/pkg/checkout"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	ErrAddToCartEmptyUser = status.Error(codes.InvalidArgument, "empty user")
	ErrAddToCartEmptySKU  = status.Error(codes.InvalidArgument, "empty sku")
)

func (i *Implementation) AddToCart(ctx context.Context, req *desc.AddToCartRequest) (*emptypb.Empty, error) {
	if req.User == 0 {
		return &emptypb.Empty{}, ErrAddToCartEmptyUser
	}
	if req.Sku == 0 {
		return &emptypb.Empty{}, ErrAddToCartEmptySKU
	}

	err := i.checkoutService.AddToCart(ctx, req.User, req.Sku, req.Count)
	if err != nil {
		return &emptypb.Empty{}, status.Error(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}
