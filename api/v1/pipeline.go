package v1

import (
	"github.com/klovercloud-ci/core/v1/api"
	"github.com/klovercloud-ci/core/v1/service"
	"github.com/labstack/echo/v4"
)

type pipelineApi struct {
	pipelineService service.Pipeline
}

func (p pipelineApi) Apply(context echo.Context) error {
	panic("implement me")
}

func (p pipelineApi) GetLog(context echo.Context) error {
	panic("implement me")
}

func NewPipelineApi(pipelineService service.Pipeline) api.Pipeline {
	return &pipelineApi{
		pipelineService: pipelineService,
	}
}
