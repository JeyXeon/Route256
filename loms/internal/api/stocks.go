package loms

import (
	"context"
	"log"
	"route256/loms/internal/converter"
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

	stocks, err := i.lomsService.GetStocks(ctx, req.Sku)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &desc.StocksResponse{
		Stocks: converter.ToStockListLomsApi(stocks),
	}, nil
}
