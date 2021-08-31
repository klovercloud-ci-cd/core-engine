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
	var logEventService service.LogEvent
	var processEventService service.ProcessEvent
	tektonClientSet,k8sClientSet:=config.GetClientSet()
	if config.Database==enums.Mongo{
		logEventService=logic.NewLogEventService(mongo.NewLogEventRepository())
		processEventService=logic.NewProcessEventService(in_memory.NewProcessEventRepository())
		tekton = logic.NewTektonService(tektonClientSet, logEventService)
		k8s=logic.NewK8sService(k8sClientSet,logEventService,processEventService,tekton)

	}
	if config.Database == enums.Inmemory{
		logEventService=logic.NewLogEventService(in_memory.NewLogEventRepository())
		processEventService=logic.NewProcessEventService(in_memory.NewProcessEventRepository())
		tekton = logic.NewTektonService(tektonClientSet,logEventService)
		k8s = logic.NewK8sService(k8sClientSet, logEventService,processEventService,tekton)
	}
	return logic.NewPipelineService(k8s,tekton,logEventService,processEventService)
}