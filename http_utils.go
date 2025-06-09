package ezutil

import (
	"math"
	"net/http"

	"github.com/itsLeonB/ezutil/internal"
)

type QueryOptions struct {
	Page  int `query:"page" binding:"required,min=1"`
	Limit int `query:"limit" binding:"required,min=1"`
}

type Pagination struct {
	TotalData   int  `json:"totalData"`
	CurrentPage int  `json:"currentPage"`
	TotalPages  int  `json:"totalPages"`
	HasNextPage bool `json:"hasNextPage"`
	HasPrevPage bool `json:"hasPrevPage"`
}

func (p *Pagination) IsZero() bool {
	return p.TotalData == 0 && p.CurrentPage == 0 && p.TotalPages == 0 && !p.HasNextPage && !p.HasPrevPage
}

type JSONResponse struct {
	Message    string     `json:"message"`
	Data       any        `json:"data,omitzero"`
	Errors     error      `json:"errors,omitempty"`
	Pagination Pagination `json:"pagination,omitzero"`
}

func NewResponse(message string) JSONResponse {
	return JSONResponse{
		Message: message,
	}
}

func NewErrorResponse(err error) any {
	return JSONResponse{
		Message: err.Error(),
		Errors:  err,
	}
}

func (jr JSONResponse) WithData(data any) JSONResponse {
	jr.Data = data
	return jr
}

func (jr JSONResponse) WithError(err error) JSONResponse {
	jr.Errors = err
	return jr
}

func (jr JSONResponse) WithPagination(queryOptions QueryOptions, totalData int) JSONResponse {
	totalPages := int(math.Ceil(float64(totalData) / float64(queryOptions.Limit)))

	jr.Pagination = Pagination{
		TotalData:   totalData,
		CurrentPage: queryOptions.Page,
		TotalPages:  totalPages,
		HasNextPage: queryOptions.Page < totalPages,
		HasPrevPage: queryOptions.Page > 1,
	}

	return jr
}

func RunServer(defaultConfigs Config, serverSetupFunc func(*Config) *http.Server) {
	configs := LoadConfig(defaultConfigs)
	srv := serverSetupFunc(configs)
	internal.ServeGracefully(srv, configs.App.Timeout)
}
