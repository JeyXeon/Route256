package loms

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
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

	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, c.urlCreateOrder, bytes.NewBuffer(rawJSON))
	if err != nil {
		return 0, errors.Wrap(err, "creating http request")
	}

	httpResponse, err := http.DefaultClient.Do(httpRequest)
	if err != nil {
		return 0, errors.Wrap(err, "calling http")
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("wrong status code: %d", httpResponse.StatusCode)
	}

	var response CreateOrderResponse
	err = json.NewDecoder(httpResponse.Body).Decode(&response)
	if err != nil {
		return 0, errors.Wrap(err, "decoding json")
	}

	return response.OrderId, nil
}
