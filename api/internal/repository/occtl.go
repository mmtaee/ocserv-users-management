package repository

import (
	"api/internal/models"
	"api/pkg/config"
	ocApi "api/pkg/oc_api"
	"context"
	"encoding/json"
)

type OcctlRepository struct {
	ocApi ocApi.OcOcctlApiRepositoryInterface
}

type OcctlRepositoryInterface interface {
	Version(c context.Context) (*models.ServerVersion, error)
	Status(c context.Context) (map[string]interface{}, error)
	OnlineUsers(c context.Context) ([]string, error)
	OnlineUsersInfo(c context.Context) ([]models.OnlineUserSession, error)
	IPBans(c context.Context) ([]models.IPBan, error)
	IRoutes(c context.Context) ([]models.Iroute, error)
	Reload(c context.Context) error
}

func NewOcctlRepository() *OcctlRepository {
	webhookApi := config.Get().WebhookApi
	return &OcctlRepository{
		ocApi: ocApi.NewOcctlApiRepository(webhookApi),
	}
}

func (o *OcctlRepository) Version(c context.Context) (*models.ServerVersion, error) {
	resp, err := o.ocApi.Version(c)
	if err != nil {
		return nil, err
	}

	var result models.ServerVersion
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (o *OcctlRepository) Status(c context.Context) (map[string]interface{}, error) {
	c = context.WithValue(c, "format", "json")
	res, err := o.ocApi.Status(c)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	err = json.Unmarshal(res, &result)
	return result, nil
}

func (o *OcctlRepository) OnlineUsers(c context.Context) ([]string, error) {
	var res []string
	users, err := o.ocApi.OnlineUsers(c)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(users, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (o *OcctlRepository) OnlineUsersInfo(c context.Context) ([]models.OnlineUserSession, error) {
	res, err := o.ocApi.OnlineUsersInfo(c)
	if err != nil {
		return nil, err
	}
	var results []models.OnlineUserSession
	err = json.Unmarshal(res, &results)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (o *OcctlRepository) IPBans(c context.Context) ([]models.IPBan, error) {
	var results []models.IPBan

	res, err := o.ocApi.IPBans(c)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(res, &results)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (o *OcctlRepository) IRoutes(c context.Context) ([]models.Iroute, error) {
	res, err := o.ocApi.IRoutes(c)
	if err != nil {
		return nil, err
	}
	var results []models.Iroute
	err = json.Unmarshal(res, &results)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (o *OcctlRepository) Reload(c context.Context) error {
	return o.ocApi.Reload(c)
}
