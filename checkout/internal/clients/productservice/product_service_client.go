package productservice

type Client struct {
	url   string
	token string

	urlGetProduct string
}

func New(url string, token string) *Client {
	return &Client{
		url:   url,
		token: token,

		urlGetProduct: url + "/get_product",
	}
}
