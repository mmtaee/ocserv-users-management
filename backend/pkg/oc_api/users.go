package oc_api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type OcUserApiRepository struct {
	url string
}

type OcUserApiRepositoryInterface interface {
	CreateUserApi(c context.Context, group, username, password string) error
	LockUserApi(c context.Context, username string) error
	UnLockUserApi(c context.Context, username string) error
	DeleteUserApi(c context.Context, username string) error
}

func NewOcUserApiRepository(url string) *OcUserApiRepository {
	return &OcUserApiRepository{url: url}
}

func (o *OcUserApiRepository) CreateUserApi(ctx context.Context, group, username, password string) error {
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
	resp, err := DoRequest(ctx, url, http.MethodPost, jsonData)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}
	return nil
}

func (o *OcUserApiRepository) LockUserApi(ctx context.Context, username string) error {
	url := fmt.Sprintf("%s/api/users/%s/lock", o.url, username)
	resp, err := DoRequest(ctx, url, http.MethodPost, nil)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}
	return nil
}

func (o *OcUserApiRepository) UnLockUserApi(ctx context.Context, username string) error {
	url := fmt.Sprintf("%s/api/users/%s/unlock", o.url, username)
	resp, err := DoRequest(ctx, url, http.MethodPost, nil)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}
	return nil
}

func (o *OcUserApiRepository) DeleteUserApi(ctx context.Context, username string) error {
	url := fmt.Sprintf("%s/api/users/%s", o.url, username)
	resp, err := DoRequest(ctx, url, http.MethodDelete, nil)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusNoContent {
		return errors.New(resp.Status)
	}
	return nil
}
