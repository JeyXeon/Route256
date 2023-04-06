package checkout

import (
	"context"
	desc "route256/checkout/pkg/checkout"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	ErrDeleteFromCartEmptyUser = status.Error(codes.InvalidArgument, "empty user")
	ErrDeleteFromCartEmptySKU  = status.Error(codes.InvalidArgument, "empty sku")
)

func (i *Implementation) DeleteFromCart(ctx context.Context, req *desc.DeleteFromCartRequest) (*emptypb.Empty, error) {
	if req.User == 0 {
		return &emptypb.Empty{}, ErrDeleteFromCartEmptyUser
	}
	if req.Sku == 0 {
		return &emptypb.Empty{}, ErrDeleteFromCartEmptySKU
	}

	err := i.checkoutService.DeleteFromCart(ctx, req.User, req.Sku, req.Count)
	if err != nil {
		return &emptypb.Empty{}, status.Error(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}
