package productservice

import "context"

type RequestProcessor interface {
	ProcessRequest(ctx context.Context, url string, requestJson []byte) ([]byte, error)
}

type Client struct {
	requestProcessor RequestProcessor
	token            string
}

const (
	urlGetProduct = "/get_product"
)

func New(requestProcessor RequestProcessor, token string) *Client {
	return &Client{
		requestProcessor: requestProcessor,
		token:            token,
	}
}
