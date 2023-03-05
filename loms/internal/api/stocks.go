package loms

import (
	"context"
	"log"
	desc "route256/loms/pkg/loms"

	"github.com/pkg/errors"
)

var (
	ErrStocksEmptySKU = errors.New("empty sku")
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
