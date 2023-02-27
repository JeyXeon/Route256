package productservice

import (
	"context"
	"encoding/json"
	"route256/checkout/internal/domain"

	"github.com/pkg/errors"
)

type GetProductRequest struct {
	Token string `json:"token"`
	SKU   uint32 `json:"sku"`
}

type GetProductResponse struct {
	Name  string `json:"name"`
	Price uint32 `json:"price"`
}

func (c *Client) GetProduct(ctx context.Context, SKU uint32) (domain.Product, error) {
	request := GetProductRequest{Token: c.token, SKU: SKU}

	rawJSON, err := json.Marshal(request)
	if err != nil {
		return domain.Product{}, errors.Wrap(err, "marshaling json")
	}

	responseJson, err := c.requestProcessor.ProcessRequest(ctx, urlGetProduct, rawJSON)
	if err != nil {
		return domain.Product{}, errors.Wrap(err, "processing request")
	}

	var response GetProductResponse
	err = json.Unmarshal(responseJson, &response)
	if err != nil {
		return domain.Product{}, errors.Wrap(err, "decoding json")
	}

	return domain.Product{SKU: SKU, Count: 1, Name: response.Name, Price: response.Price}, nil
}
