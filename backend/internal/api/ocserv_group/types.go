package ocserv_group

import (
	"ocserv-bakend/internal/models"
	"ocserv-bakend/pkg/request"
)

type OcservGroupsResponse struct {
	Meta   request.Meta          `json:"meta" validate:"required"`
	Result *[]models.OcservGroup `json:"result" validate:"omitempty"`
}
