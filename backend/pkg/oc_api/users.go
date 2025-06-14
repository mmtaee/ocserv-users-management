package oc_api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type OcUserApiRepository struct {
	url string
}

type OcUserApiRepositoryInterface interface {
	CreateApi(group, username, password string) error
	LockApi(username string) error
	UnLockApi(username string) error
	DeleteApi(username string) error
}

func NewOcUserApiRepository(url string) *OcUserApiRepository {
	return &OcUserApiRepository{url: url}
}

func (o *OcUserApiRepository) CreateApi(group, username, password string) error {
	url := o.url + "/api/users"
	type Body struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Group    string `json:"group"`
	}
	data := Body{
		Username: username,
		Password: password,
		Group:    group,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	resp, err := DoRequest(url, http.MethodPost, jsonData)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusCreated {
		return errors.New(resp.Status)
	}
	return nil
}

func (o *OcUserApiRepository) LockApi(username string) error {
	url := fmt.Sprintf("%s/api/users/%s/lock", o.url, username)
	resp, err := DoRequest(url, http.MethodPost, nil)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusCreated {
		return errors.New(resp.Status)
	}
	return nil
}

func (o *OcUserApiRepository) UnLockApi(username string) error {
	url := fmt.Sprintf("%s/api/users/%s/unlock", o.url, username)
	resp, err := DoRequest(url, http.MethodPost, nil)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusCreated {
		return errors.New(resp.Status)
	}
	return nil
}

func (o *OcUserApiRepository) DeleteApi(username string) error {
	url := fmt.Sprintf("%s/api/users/%s", o.url, username)
	resp, err := DoRequest(url, http.MethodDelete, nil)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusNoContent {
		return errors.New(resp.Status)
	}
	return nil
}
