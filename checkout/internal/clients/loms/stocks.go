package loms

import (
	"context"
	"route256/checkout/internal/model"
	lomsapi "route256/checkout/pkg/loms"
)

func (c *client) Stocks(ctx context.Context, sku uint32) ([]*model.Stock, error) {
	req := &lomsapi.StocksRequest{
		Sku: sku,
	}

	res, err := c.lomsClient.Stocks(ctx, req)
	if err != nil {
		return nil, err
	}

	stocks := make([]*model.Stock, 0, len(res.GetStocks()))
	for _, stock := range res.GetStocks() {
		stocks = append(stocks, &model.Stock{
			WarehouseID: stock.WarehouseId,
			Count:       stock.Count,
		})
	}

	return stocks, nil

}
