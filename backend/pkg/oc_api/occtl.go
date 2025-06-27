package oc_api

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type OcctlApiRepository struct {
	url string
}

type OcOcctlApiRepositoryInterface interface {
	Stats(ctx context.Context) (string, error)
	OnlineUsers(ctx context.Context) (*[]string, error)
}

func NewOcctlApiRepository(url string) *OcctlApiRepository {
	return &OcctlApiRepository{url: url}
}

func (o *OcctlApiRepository) Stats(ctx context.Context) (string, error) {
	url := o.url + "/api/occtl/status"
	resp, err := DoRequest(ctx, url, http.MethodGet, nil)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", errors.New(resp.Status)
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(respBody), nil
}

func (o *OcctlApiRepository) OnlineUsers(ctx context.Context) (*[]string, error) {
	url := o.url + "/api/occtl/online-users"
	resp, err := DoRequest(ctx, url, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result struct {
		Users []string `json:"users"`
	}

	err = json.Unmarshal(respBody, &result)
	if err != nil {
		return nil, err
	}

	return &result.Users, nil
}
