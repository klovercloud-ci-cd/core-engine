package v1

import (
	"github.com/klovercloud-ci/api/common"
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
		return common.GenerateErrorResponse(context,nil,err.Error())
	}
	url:=context.QueryParam("url")
	revision:=context.QueryParam("revision")

	data.ApiVersion = enums.Api_version
	data.ProcessId = guuid.New().String()
	error := p.pipelineService.Apply(url,revision, data)

	if error != nil{
		log.Println("Input Error:", err.Error())
		return common.GenerateErrorResponse(context,err.Error(),"Failed to trigger pipeline!")
	}
	return common.GenerateSuccessResponse(context,data.ProcessId,nil,"Pipeline successfully triggered!")
}

func NewPipelineApi(pipelineService service.Pipeline) api.Pipeline {
	return &pipelineApi{
		pipelineService: pipelineService,
	}
}
