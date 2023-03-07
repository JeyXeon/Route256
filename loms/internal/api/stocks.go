package loms

import (
	"context"
	"log"
	desc "route256/loms/pkg/loms"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrStocksEmptySKU = status.Error(codes.InvalidArgument, "empty sku")
)

func (i *Implementation) Stocks(ctx context.Context, req *desc.StocksRequest) (*desc.StocksResponse, error) {
	if req.Sku == 0 {
		return nil, ErrStocksEmptySKU
	}

	log.Printf("stocks: %+v", req)

	return &desc.StocksResponse{
		Stocks: []*desc.Stock{
			{
				WarehouseId: 123,
				Count:       5,
			},
		},
	}, nil
}
