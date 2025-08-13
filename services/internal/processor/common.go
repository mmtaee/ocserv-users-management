package processor

import (
	"services/pkg/config"
	ocApi "services/pkg/oc_api"
)

var (
	occtlApi  ocApi.OcOcctlApiRepositoryInterface
	ocUserApi ocApi.OcUserApiRepositoryInterface
)

func Init() {
	cfg := config.Get()

	occtlApi = ocApi.NewOcctlApiRepository(cfg.APIURLService)
	ocUserApi = ocApi.NewOcUserApiRepository(cfg.APIURLService)
}
