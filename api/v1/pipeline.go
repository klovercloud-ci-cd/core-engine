package v1

import (
	v1 "github.com/klovercloud-ci/core/v1"
	"github.com/klovercloud-ci/core/v1/api"
	"github.com/klovercloud-ci/core/v1/service"
	"github.com/klovercloud-ci/enums"
	"github.com/labstack/echo/v4"
	guuid "github.com/google/uuid"
	"log"
)

type pipelineApi struct {
	pipelineService service.Pipeline
}

func (p pipelineApi) GetLog(context echo.Context)error {
	panic("implement me")
}

func (p pipelineApi) GetEvents(context echo.Context) error {
	panic("implement me")
}

func (p pipelineApi) Apply(context echo.Context) error {
	data:=v1.Pipeline{}
	err := context.Bind(&data)
	if  err != nil{
		log.Println("Input Error:", err.Error())
		return err
	}
	data.ApiVersion = enums.Api_version
	pId := guuid.New()
	data.ProcessId = pId.String()
	error := p.pipelineService.Apply("www.example.com","www.revision.com", data)

	if error != nil{
		return error
	}
	return nil
}

func NewPipelineApi(pipelineService service.Pipeline) api.Pipeline {
	return &pipelineApi{
		pipelineService: pipelineService,
	}
}
