package oc_api

import (
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
	CreateGroupApi(name string, config *models.OcservGroupConfig) error
	DeleteGroupApi(name string) error
}

func NewOcGroupApiRepository(url string) *OcGroupApiRepository {
	return &OcGroupApiRepository{url: url}
}

func (o *OcGroupApiRepository) CreateGroupApi(name string, config *models.OcservGroupConfig) error {
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

	resp, err := DoRequest(url, http.MethodPost, jsonData)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}
	return nil
}

func (o *OcGroupApiRepository) DeleteGroupApi(name string) error {
	url := fmt.Sprintf("%s/api/groups/%s", o.url, name)
	resp, err := DoRequest(url, http.MethodDelete, nil)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusNoContent {
		return errors.New(resp.Status)
	}
	return nil
}
