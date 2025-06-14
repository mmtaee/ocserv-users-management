package ocserv_group

import (
	"ocserv-bakend/internal/models"
	"ocserv-bakend/pkg/request"
)

type CreateOcservGroupData struct {
	Name   string                    `json:"name" validate:"required"`
	Config *models.OcservGroupConfig `json:"config" validate:"required"`
}

type UpdateOcservGroupData struct {
	Config *models.OcservGroupConfig `json:"config" validate:"required"`
}

type OcservGroupsResponse struct {
	Meta   request.Meta          `json:"meta" validate:"required"`
	Result *[]models.OcservGroup `json:"result" validate:"omitempty"`
}
