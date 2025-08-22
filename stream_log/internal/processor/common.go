package processor

import (
	"stream_log/pkg/config"
	ocApi "stream_log/pkg/oc_api"
)

var (
	occtlApi  ocApi.OcOcctlApiRepositoryInterface
	ocUserApi ocApi.OcUserApiRepositoryInterface
)

func Init() {
	cfg := config.Get()

	occtlApi = ocApi.NewOcctlApiRepository(cfg.WebhookApi)
	ocUserApi = ocApi.NewOcUserApiRepository(cfg.WebhookApi)
}
