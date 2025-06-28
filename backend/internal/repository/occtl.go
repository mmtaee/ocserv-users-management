package repository

import (
	"context"
	"encoding/json"
	"ocserv-bakend/internal/models"
	"ocserv-bakend/pkg/config"
	ocApi "ocserv-bakend/pkg/oc_api"
)

type OcctlRepository struct {
	ocApi ocApi.OcOcctlApiRepositoryInterface
}

type OcctlRepositoryInterface interface {
	Stats(c context.Context) (string, error)
	OnlineUsers(c context.Context) (*[]string, error)
	OnlineUsersInfo(c context.Context) (*[]models.OnlineUserSession, error)
	IPBans(c context.Context) (*[]models.IPBan, error)
	IRoutes(c context.Context) (*[]models.Iroute, error)
}

func NewOcctlRepository() *OcctlRepository {
	apiURLService := config.Get().APIURLService
	return &OcctlRepository{
		ocApi: ocApi.NewOcctlApiRepository(apiURLService),
	}
}

func (o *OcctlRepository) Stats(c context.Context) (string, error) {
	return o.ocApi.Stats(c)
}

func (o *OcctlRepository) OnlineUsers(c context.Context) (*[]string, error) {
	return o.ocApi.OnlineUsers(c)
}

func (o *OcctlRepository) OnlineUsersInfo(c context.Context) (*[]models.OnlineUserSession, error) {
	res, err := o.ocApi.OnlineUsersInfo(c)
	if err != nil {
		return nil, err
	}
	var results []models.OnlineUserSession
	err = json.Unmarshal(*res, &results)
	if err != nil {
		return nil, err
	}
	return &results, nil
}

func (o *OcctlRepository) IPBans(c context.Context) (*[]models.IPBan, error) {
	var results []models.IPBan

	res, err := o.ocApi.IPBans(c)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(*res, &results)
	if err != nil {
		return nil, err
	}
	return &results, nil
}

func (o *OcctlRepository) IRoutes(c context.Context) (*[]models.Iroute, error) {
	res, err := o.ocApi.IRoutes(c)
	if err != nil {
		return nil, err
	}
	var results []models.Iroute
	err = json.Unmarshal(*res, &results)
	if err != nil {
		return nil, err
	}
	return &results, nil
}
