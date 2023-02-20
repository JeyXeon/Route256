package productservice

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"route256/checkout/internal/domain"
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

	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, c.urlGetProduct, bytes.NewBuffer(rawJSON))
	if err != nil {
		return domain.Product{}, errors.Wrap(err, "creating http request")
	}

	httpResponse, err := http.DefaultClient.Do(httpRequest)
	if err != nil {
		return domain.Product{}, errors.Wrap(err, "calling http")
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode != http.StatusOK {
		return domain.Product{}, fmt.Errorf("wrong status code: %d", httpResponse.StatusCode)
	}

	var response GetProductResponse
	err = json.NewDecoder(httpResponse.Body).Decode(&response)
	if err != nil {
		return domain.Product{}, errors.Wrap(err, "decoding json")
	}

	return domain.Product{SKU: SKU, Count: 1, Name: response.Name, Price: response.Price}, nil
}
