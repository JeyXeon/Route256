package checkout

import (
	"context"
	desc "route256/checkout/pkg/checkout"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrPurchaseEmptyUser = status.Error(codes.InvalidArgument, "empty user")
)

func (i *Implementation) Purchase(ctx context.Context, req *desc.PurchaseRequest) (*desc.PurchaseResponse, error) {
	if req.User == 0 {
		return nil, ErrPurchaseEmptyUser
	}

	order, err := i.checkoutService.Purchase(ctx, req.User)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &desc.PurchaseResponse{
		OrderID: order,
	}, nil
}
