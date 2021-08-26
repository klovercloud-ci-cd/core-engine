package dependency

import (
	"github.com/klovercloud-ci/core/v1/logic"
	"github.com/klovercloud-ci/core/v1/service"
)

func GetPipelineService() service.Pipeline{
	var k8s service.K8s
	var tekon service.Tekton
	return logic.NewPipelineService(k8s,tekon)
}