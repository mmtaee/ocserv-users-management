package oc_api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"ocserv-bakend/internal/models"
)

type OcGroupApiRepository struct {
	url string
}

type OcGroupApiRepositoryInterface interface {
	CreateGroupApi(ctx context.Context, name string, config *models.OcservGroupConfig) error
	DeleteGroupApi(ctx context.Context, name string) error
	GetDefaultsGroup(ctx context.Context) (map[string]interface{}, error)
	UpdateDefaultGroup(ctx context.Context, config *models.OcservGroupConfig) error
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
		bodyBytes, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			return fmt.Errorf("request failed with status %d and could not read body: %w", resp.StatusCode, readErr)
		}
		return fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}
	return nil
}

func (o *OcGroupApiRepository) DeleteGroupApi(ctx context.Context, name string) error {
	url := fmt.Sprintf("%s/api/groups/%s", o.url, name)
	resp, err := DoRequest(ctx, url, http.MethodDelete, nil)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}
	return nil
}

func (o *OcGroupApiRepository) GetDefaultsGroup(ctx context.Context) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/api/groups/defaults", o.url)
	resp, err := DoRequest(ctx, url, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}
	result := map[string]interface{}{}
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

func (o *OcGroupApiRepository) UpdateDefaultGroup(ctx context.Context, config *models.OcservGroupConfig) error {
	url := fmt.Sprintf("%s/api/groups/defaults", o.url)

	jsonData, err := json.Marshal(config)
	if err != nil {
		return err
	}
	
	resp, err := DoRequest(ctx, url, http.MethodPatch, jsonData)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		bodyBytes, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			return fmt.Errorf("request failed with status %d and could not read body: %w", resp.StatusCode, readErr)
		}
		return fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}
	return nil
}
