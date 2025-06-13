package request

import "github.com/labstack/echo/v4"

// Pagination defines pagination query parameters for the API.
// @Param page query int false "Page number, starting from 1" minimum(1)
// @Param page_size query int false "Number of items per page" minimum(1) maximum(100)
// @Param order query string false "Field to order by"
// @Param sort query string false "Sort order, either ASC or DESC" Enums(ASC, DESC)
// @Description Pagination parameters
type Pagination struct {
	Page     int    `json:"page" query:"page" validate:"omitempty,min=1"`
	PageSize int    `json:"page_size" query:"page_size" validate:"omitempty,min=1,max=100"`
	Order    string `json:"order" query:"order" validate:"omitempty"`
	Sort     string `json:"sort" query:"sort" validate:"omitempty,oneof=DESC ASC"`
}

func (r *Request) Pagination(c echo.Context) *Pagination {
	var pagination Pagination

	if err := c.Bind(&pagination); err != nil {
		return &Pagination{
			Page:     1,
			PageSize: 50,
			Order:    "id",
			Sort:     "ASC",
		}
	}

	if pagination.Page == 0 {
		pagination.Page = 1
	}
	if pagination.PageSize == 0 {
		pagination.PageSize = 50
	}
	if pagination.Order == "" {
		pagination.Order = "id"
	}
	if pagination.Sort == "" {
		pagination.Sort = "ASC"
	}

	return &pagination
}

type Meta struct {
	Page         int   `json:"page" validate:"required"`
	PageSize     int   `json:"page_size" validate:"required"`
	TotalRecords int64 `json:"total_records" validate:"required"`
}
