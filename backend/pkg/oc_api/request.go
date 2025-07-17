package oc_api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
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

func (o *OcctlApiRepository) doRequestBytes(ctx context.Context, endpoint string, method string) ([]byte, error) {
	resp, err := DoRequest(ctx, endpoint, method, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)

		var errResp struct {
			Message string `json:"message"`
		}

		if json.Unmarshal(bodyBytes, &errResp) == nil && errResp.Message != "" {
			return nil, errors.New(errResp.Message)
		}

		return nil, errors.New(strings.TrimSpace(string(bodyBytes)))
	}

	return io.ReadAll(resp.Body)
}
