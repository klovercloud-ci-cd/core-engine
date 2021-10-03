package dependency

import (
	"github.com/klovercloud-ci/config"
	"github.com/klovercloud-ci/core/v1/logic"
	"github.com/klovercloud-ci/core/v1/service"
	"github.com/klovercloud-ci/enums"
	in_memory "github.com/klovercloud-ci/repository/v1/in-memory"
	"github.com/klovercloud-ci/repository/v1/mongo"
)

func GetObserverServices()[]service.Observer {
	var observers [] service.Observer
	var processEventService service.ProcessEvent
	var eventStoreLogEventService service.EventStoreLogEvent
	var eventStoreProcessEvent service.EventStoreProcessEvent
	var processLifeCycleEvent service.ProcessLifeCycleEvent
	if config.UseLocalEventStore{
		if config.Database==enums.Mongo{
			processEventService=logic.NewProcessEventService(in_memory.NewProcessEventRepository())
			observers= append(observers, processEventService)
			observers= append(observers, processLifeCycleEvent)
		}
		if config.Database == enums.Inmemory{
			processEventService=logic.NewProcessEventService(in_memory.NewProcessEventRepository())
			eventStoreLogEventService=logic.NewV1EventStoreLogEventService(logic.NewHttpPublisherService())
			eventStoreProcessEvent=logic.NewV1EventStoreProcessEventService(logic.NewHttpPublisherService())
			observers= append(observers, processEventService)
			observers= append(observers, eventStoreLogEventService)
			observers= append(observers, eventStoreProcessEvent)
		}
	}else{
		eventStoreLogEventService=logic.NewV1EventStoreLogEventService(logic.NewHttpPublisherService())
		eventStoreProcessEvent=logic.NewV1EventStoreProcessEventService(logic.NewHttpPublisherService())
		observers= append(observers, eventStoreLogEventService)
		observers= append(observers, eventStoreProcessEvent)
	}
	return observers
}
func GetPipelineService() service.Pipeline {
	var k8s service.K8s
	var tekton service.Tekton
	var logEventService service.LogEvent
	var processEventService service.ProcessEvent
	var eventStoreLogEventService service.EventStoreLogEvent
	var eventStoreProcessEvent service.EventStoreProcessEvent
	var processLifeCycleEvent service.ProcessLifeCycleEvent
	var eventStoreProcessLifeCycleEvent service.EventStoreProcessLifeCycleEvent
	var observers [] service.Observer
	tektonClientSet,k8sClientSet:=config.GetClientSet()

	if config.UseLocalEventStore{
		if config.Database==enums.Mongo{
			logEventService=logic.NewLogEventService(mongo.NewLogEventRepository(3000))
			processEventService=logic.NewProcessEventService(in_memory.NewProcessEventRepository())
			observers= append(observers, logEventService)
			observers= append(observers, processEventService)
			observers= append(observers, processLifeCycleEvent)
			tekton = logic.NewTektonService(tektonClientSet)
			k8s=logic.NewK8sService(k8sClientSet,tekton,observers)

		}
		if config.Database == enums.Inmemory{
			logEventService=logic.NewLogEventService(in_memory.NewLogEventRepository())
			processEventService=logic.NewProcessEventService(in_memory.NewProcessEventRepository())
			eventStoreLogEventService=logic.NewV1EventStoreLogEventService(logic.NewHttpPublisherService())
			eventStoreProcessEvent=logic.NewV1EventStoreProcessEventService(logic.NewHttpPublisherService())
			observers= append(observers, logEventService)
			observers= append(observers, processEventService)
			observers= append(observers, eventStoreLogEventService)
			observers= append(observers, eventStoreProcessEvent)
			tekton = logic.NewTektonService(tektonClientSet)
			k8s=logic.NewK8sService(k8sClientSet,tekton,observers)
		}
	}else{
		eventStoreLogEventService=logic.NewV1EventStoreLogEventService(logic.NewHttpPublisherService())
		eventStoreProcessEvent=logic.NewV1EventStoreProcessEventService(logic.NewHttpPublisherService())
		eventStoreProcessLifeCycleEvent=logic.NewEventStoreProcessLifeCycleService(logic.NewHttpPublisherService())
		observers= append(observers, eventStoreLogEventService)
		observers= append(observers, eventStoreProcessEvent)
		observers= append(observers, eventStoreProcessLifeCycleEvent)
		tekton = logic.NewTektonService(tektonClientSet)
		k8s=logic.NewK8sService(k8sClientSet,tekton,observers)
	}

	return logic.NewPipelineService(k8s,tekton,logEventService,processEventService,observers)
}