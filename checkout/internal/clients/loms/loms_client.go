package loms

type Client struct {
	url string

	urlStocks      string
	urlCreateOrder string
}

func New(url string) *Client {
	return &Client{
		url: url,

		urlStocks:      url + "/stocks",
		urlCreateOrder: url + "/createOrder",
	}
}
