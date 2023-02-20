package loms

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"route256/checkout/internal/domain"
)

type CreateOrderRequest struct {
	User  int64              `json:"user"`
	Items []domain.OrderItem `json:"items"`
}

type CreateOrderResponse struct {
	OrderId int64 `json:"orderId"`
}

func (c *Client) CreateOrder(ctx context.Context, user int64, items []domain.OrderItem) (int64, error) {
	request := CreateOrderRequest{User: user, Items: items}

	rawJSON, err := json.Marshal(request)
	if err != nil {
		return 0, errors.Wrap(err, "marshaling json")
	}

	responseJson, err := c.lomsRequestProcessor.ProcessRequest(ctx, urlCreateOrder, rawJSON)
	if err != nil {
		return 0, errors.Wrap(err, "processing request")
	}

	var response CreateOrderResponse
	err = json.Unmarshal(responseJson, &response)
	if err != nil {
		return 0, errors.Wrap(err, "decoding json")
	}

	return response.OrderId, nil
}
