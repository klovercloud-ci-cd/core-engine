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

	tektonClientSet,k8sClientSet:=config.GetClientSet()
	if config.Database==enums.Mongo{
		tekton = logic.NewTektonService(tektonClientSet, mongo.NewLogEventRepository())
		k8s=logic.NewK8sService(k8sClientSet,logic.NewLogEventService(mongo.NewLogEventRepository()),logic.NewProcessEventService(in_memory.NewProcessEventRepository()),tekton)

	}
	if config.Database == enums.Inmemory{
		tekton = logic.NewTektonService(tektonClientSet,in_memory.NewLogEventRepository())
		k8s = logic.NewK8sService(k8sClientSet, logic.NewLogEventService(in_memory.NewLogEventRepository()),logic.NewProcessEventService(in_memory.NewProcessEventRepository()),tekton)
	}
	return logic.NewPipelineService(k8s,tekton)
}