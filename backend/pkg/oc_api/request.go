package oc_api

import (
	"bytes"
	"context"
	"net/http"
	"time"
)

func DoRequest(ctx context.Context, url string, method string, body []byte) (*http.Response, error) {
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	return client.Do(req)
}
