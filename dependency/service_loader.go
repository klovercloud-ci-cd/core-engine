package dependency

import (
	"github.com/klovercloud-ci/config"
	"github.com/klovercloud-ci/core/v1/logic"
	"github.com/klovercloud-ci/core/v1/service"
	"github.com/klovercloud-ci/enums"
	in_memory "github.com/klovercloud-ci/repository/v1/inmemory"
	"github.com/klovercloud-ci/repository/v1/mongo"
)

// GetV1ObserverServices returns Observer services
func GetV1ObserverServices() []service.Observer {
	var observers []service.Observer
	var processEventService service.ProcessEvent
	var eventStoreLogEventService service.Observer
	var eventStoreProcessEvent service.Observer
	var processLifeCycleEvent service.ProcessLifeCycleEvent
	if config.UseLocalEventStore {
		if config.Database == enums.MONGO {
			processEventService = logic.NewProcessEventService(in_memory.NewProcessEventRepository())
			observers = append(observers, processEventService)
			observers = append(observers, processLifeCycleEvent)
		}
		if config.Database == enums.INMEMORY {
			processEventService = logic.NewProcessEventService(in_memory.NewProcessEventRepository())
			eventStoreLogEventService = logic.NewEventStoreLogEventService(logic.NewHttpClientService())
			eventStoreProcessEvent = logic.NewEventStoreProcessEventService(logic.NewHttpClientService())
			observers = append(observers, processEventService)
			observers = append(observers, eventStoreLogEventService)
			observers = append(observers, eventStoreProcessEvent)
		}
	} else {
		eventStoreLogEventService = logic.NewEventStoreLogEventService(logic.NewHttpClientService())
		eventStoreProcessEvent = logic.NewEventStoreProcessEventService(logic.NewHttpClientService())
		observers = append(observers, eventStoreLogEventService)
		observers = append(observers, eventStoreProcessEvent)
	}
	return observers
}

// GetV1PipelineService returns Pipeline services
func GetV1PipelineService() service.Pipeline {
	var k8s service.K8s
	var tekton service.Tekton
	var logEventService service.LogEvent
	var processEventService service.ProcessEvent
	var eventStoreLogEventService service.Observer
	var eventStoreProcessEvent service.Observer
	var processLifeCycleEvent service.ProcessLifeCycleEvent
	var eventStoreProcessLifeCycleEvent service.ProcessLifeCycleEvent
	var observers []service.Observer
	tektonClientSet, k8sClientSet := config.GetClientSet()

	if config.UseLocalEventStore {
		if config.Database == enums.MONGO {
			logEventService = logic.NewLogEventService(mongo.NewLogEventRepository(3000))
			processEventService = logic.NewProcessEventService(in_memory.NewProcessEventRepository())
			processLifeCycleEvent = logic.NewProcessLifeCycleEventService(mongo.NewProcessLifeCycleRepository(3000))
			observers = append(observers, logEventService)
			observers = append(observers, processEventService)
			observers = append(observers, processLifeCycleEvent)
			tekton = logic.NewTektonService(tektonClientSet)
			k8s = logic.NewK8sService(k8sClientSet, tekton, observers)

		}
		if config.Database == enums.INMEMORY {
			logEventService = logic.NewLogEventService(in_memory.NewLogEventRepository())
			processEventService = logic.NewProcessEventService(in_memory.NewProcessEventRepository())
			eventStoreLogEventService = logic.NewEventStoreLogEventService(logic.NewHttpClientService())
			eventStoreProcessEvent = logic.NewEventStoreProcessEventService(logic.NewHttpClientService())
			processLifeCycleEvent = logic.NewProcessLifeCycleEventService(in_memory.NewProcessLifeCycleRepository())
			observers = append(observers, logEventService)
			observers = append(observers, processEventService)
			observers = append(observers, eventStoreLogEventService)
			observers = append(observers, eventStoreProcessEvent)
			observers = append(observers, processLifeCycleEvent)
			tekton = logic.NewTektonService(tektonClientSet)
			k8s = logic.NewK8sService(k8sClientSet, tekton, observers)
		}
	} else {
		eventStoreLogEventService = logic.NewEventStoreLogEventService(logic.NewHttpClientService())
		eventStoreProcessEvent = logic.NewEventStoreProcessEventService(logic.NewHttpClientService())
		eventStoreProcessLifeCycleEvent = logic.NewEventStoreProcessLifeCycleService(logic.NewHttpClientService())
		observers = append(observers, eventStoreLogEventService)
		observers = append(observers, eventStoreProcessEvent)
		observers = append(observers, eventStoreProcessLifeCycleEvent)
		tekton = logic.NewTektonService(tektonClientSet)
		k8s = logic.NewK8sService(k8sClientSet, tekton, observers)
	}

	return logic.NewPipelineService(k8s, tekton, logEventService, processEventService, observers, eventStoreProcessLifeCycleEvent)
}

// GetV1JwtService returns Jwt services
func GetV1JwtService() service.Jwt {
	return logic.NewJwtService()
}
