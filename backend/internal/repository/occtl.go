package repository

import (
	"context"
	"ocserv-bakend/pkg/config"
	ocApi "ocserv-bakend/pkg/oc_api"
)

type OcctlRepository struct {
	ocApi ocApi.OcOcctlApiRepositoryInterface
}

type OcctlRepositoryInterface interface {
	Stats(c context.Context) (string, error)
	OnlineUsers(c context.Context) (*[]string, error)
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
