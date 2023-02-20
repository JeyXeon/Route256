package loms

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"route256/checkout/internal/domain"
)

type StocksRequest struct {
	SKU uint32 `json:"sku"`
}

type StocksItem struct {
	WarehouseID int64  `json:"warehouseID"`
	Count       uint64 `json:"count"`
}

type StocksResponse struct {
	Stocks []StocksItem `json:"stocks"`
}

func (c *Client) Stocks(ctx context.Context, sku uint32) ([]domain.Stock, error) {
	request := StocksRequest{SKU: sku}

	rawJSON, err := json.Marshal(request)
	if err != nil {
		return nil, errors.Wrap(err, "marshaling json")
	}

	responseJson, err := c.lomsRequestProcessor.ProcessRequest(ctx, urlStocks, rawJSON)
	if err != nil {
		return nil, errors.Wrap(err, "processing request")
	}

	var response StocksResponse
	err = json.Unmarshal(responseJson, &response)
	if err != nil {
		return nil, errors.Wrap(err, "decoding json")
	}

	stocks := make([]domain.Stock, 0, len(response.Stocks))
	for _, stock := range response.Stocks {
		stocks = append(stocks, domain.Stock{
			WarehouseID: stock.WarehouseID,
			Count:       stock.Count,
		})
	}

	return stocks, nil
}
