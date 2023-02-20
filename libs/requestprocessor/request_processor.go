package requestprocessor

import (
	"bytes"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"net/http"
)

type RequestProcessor struct {
	rootUrl string
}

func New(rootUrl string) *RequestProcessor {
	return &RequestProcessor{
		rootUrl: rootUrl,
	}
}

func (r *RequestProcessor) ProcessRequest(ctx context.Context, url string, request []byte) ([]byte, error) {
	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, r.rootUrl+url, bytes.NewBuffer(request))
	if err != nil {
		return nil, errors.Wrap(err, "creating http request")
	}

	httpResponse, err := http.DefaultClient.Do(httpRequest)
	if err != nil {
		return nil, errors.Wrap(err, "calling http")
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("wrong status code: %d", httpResponse.StatusCode)
	}

	response, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, errors.Wrap(err, "decoding response")
	}

	return response, nil
}
