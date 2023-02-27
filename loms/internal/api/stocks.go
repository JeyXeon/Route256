package loms

import (
	"context"
	"log"
	desc "route256/loms/pkg"
)

func (i *Implementation) Stocks(ctx context.Context, req *desc.StocksRequest) (*desc.StocksResponse, error) {
	if req.Sku == 0 {
		return nil, ErrEmptySKU
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
