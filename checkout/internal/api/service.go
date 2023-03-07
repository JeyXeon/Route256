package checkout

import (
	"route256/checkout/internal/service"
	desc "route256/checkout/pkg/checkout"
)

type Implementation struct {
	desc.UnimplementedCheckoutServer

	checkoutService *service.Service
}

func NewCheckout(checkoutService *service.Service) *Implementation {
	return &Implementation{
		desc.UnimplementedCheckoutServer{},
		checkoutService,
	}
}
