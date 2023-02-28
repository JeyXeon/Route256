package checkout

import (
	"context"
	"log"
	desc "route256/checkout/pkg/checkout"
)

func (i *Implementation) Purchase(ctx context.Context, req *desc.PurchaseRequest) (*desc.PurchaseResponse, error) {
	if req.User == 0 {
		return nil, ErrEmptyUser
	}

	log.Printf("purchase: %+v", req)

	order, err := i.checkoutService.Purchase(ctx, req.User)
	if err != nil {
		return nil, err
	}

	return &desc.PurchaseResponse{
		OrderID: order,
	}, nil
}
