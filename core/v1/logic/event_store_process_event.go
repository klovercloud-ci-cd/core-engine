package logic

import (
	"encoding/json"
	"github.com/klovercloud-ci-cd/core-engine/config"
	"github.com/klovercloud-ci-cd/core-engine/core/v1"
	"github.com/klovercloud-ci-cd/core-engine/core/v1/service"
	"log"
)

type eventStoreProcessEventService struct {
	httpPublisher service.HttpClient
}

func (e eventStoreProcessEventService) Listen(subject v1.Subject) {
	if subject.EventData != nil {
		event := v1.ProcessEvent{
			ProcessId: subject.Pipeline.ProcessId,
			Data:      subject.EventData,
		}
		header := make(map[string]string)
		header["Content-Type"] = "application/json"
		header["token"] = config.Token
		b, err := json.Marshal(event)
		if err != nil {
			log.Println(err.Error())
			return
		}
		e.httpPublisher.Post(config.EventStoreUrl+"/processes_events", header, b)
	}
}

// NewEventStoreProcessEventService returns Observer type service
func NewEventStoreProcessEventService(httpPublisher service.HttpClient) service.Observer {
	return &eventStoreProcessEventService{
		httpPublisher: httpPublisher,
	}
}
