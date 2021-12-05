package http

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/fwojciec/asr"
	"github.com/hashicorp/go-retryablehttp"
)

type getter struct {
	client *retryablehttp.Client
}

func (g *getter) Get(ctx context.Context, url string) (io.ReadCloser, error) {
	req, err := retryablehttp.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	res, err := g.client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}
	return res.Body, nil
}

func standardClient() *retryablehttp.Client {
	client := retryablehttp.NewClient()
	client.HTTPClient.Timeout = 5 * time.Second
	return client
}

func NewGetter() asr.Getter {
	return &getter{client: standardClient()}
}
