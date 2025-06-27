package oc_api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"ocserv-bakend/internal/models"
)

type OcGroupApiRepository struct {
	url string
}

type OcGroupApiRepositoryInterface interface {
	CreateGroupApi(ctx context.Context, name string, config *models.OcservGroupConfig) error
	DeleteGroupApi(ctx context.Context, name string) error
}

func NewOcGroupApiRepository(url string) *OcGroupApiRepository {
	return &OcGroupApiRepository{url: url}
}

func (o *OcGroupApiRepository) CreateGroupApi(ctx context.Context, name string, config *models.OcservGroupConfig) error {
	url := o.url + "/api/groups"

	type Body struct {
		Name   string                    `json:"name"`
		Config *models.OcservGroupConfig `json:"config"`
	}

	body := Body{
		Name:   name,
		Config: config,
	}

	jsonData, err := json.Marshal(body)
	if err != nil {
		return err
	}

	resp, err := DoRequest(ctx, url, http.MethodPost, jsonData)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}
	return nil
}

func (o *OcGroupApiRepository) DeleteGroupApi(ctx context.Context, name string) error {
	url := fmt.Sprintf("%s/api/groups/%s", o.url, name)
	resp, err := DoRequest(ctx, url, http.MethodDelete, nil)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusNoContent {
		return errors.New(resp.Status)
	}
	return nil
}
