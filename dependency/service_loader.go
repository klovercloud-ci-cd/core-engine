package dependency

import (
	"github.com/klovercloud-ci/core/v1/logic"
	"github.com/klovercloud-ci/core/v1/service"
)

func GetPipelineService() service.Pipeline{
	k8s:=logic.NewK8sService(nil,nil,nil)
	tekon:=logic.NewTektonService(nil,nil)
	return logic.NewPipelineService(k8s,tekon)
}