package loms

import "context"

type LomsRequestProcessor interface {
	ProcessRequest(ctx context.Context, url string, requestJson []byte) ([]byte, error)
}

type Client struct {
	lomsRequestProcessor LomsRequestProcessor
}

const (
	urlStocks      = "/stocks"
	urlCreateOrder = "/createOrder"
)

func New(lomsRequestProcessor LomsRequestProcessor) *Client {
	return &Client{
		lomsRequestProcessor: lomsRequestProcessor,
	}
}
