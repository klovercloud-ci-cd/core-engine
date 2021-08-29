package dependency

import (
	"github.com/klovercloud-ci/config"
	"github.com/klovercloud-ci/core/v1/logic"
	"github.com/klovercloud-ci/core/v1/service"
	"github.com/klovercloud-ci/enums"
	in_memory "github.com/klovercloud-ci/repository/v1/in-memory"
	"github.com/klovercloud-ci/repository/v1/mongo"
)

func GetPipelineService() service.Pipeline{
	var k8s service.K8s
	var tekton service.Tekton
	if config.Database==enums.Mongo{
		k8s=logic.NewK8sService(nil,mongo.NewLogEventRepository(),in_memory.NewProcessEventRepository(),nil)
		tekton = logic.NewTektonService(nil, mongo.NewLogEventRepository())
	}
	if config.Database == enums.Inmemory{
		k8s = logic.NewK8sService(nil, in_memory.NewLogEventRepository(),in_memory.NewProcessEventRepository(),nil)
		tekton = logic.NewTektonService(nil,in_memory.NewLogEventRepository())
	}
	return logic.NewPipelineService(k8s,tekton)
}