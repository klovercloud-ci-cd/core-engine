package logic

import (
	"encoding/json"
	"github.com/klovercloud-ci-cd/klovercloud-ci-core/config"
	v1 "github.com/klovercloud-ci-cd/klovercloud-ci-core/core/v1"
	"github.com/klovercloud-ci-cd/klovercloud-ci-core/core/v1/service"
	"log"
)

type eventStoreEventService struct {
	httpPublisher service.HttpClient
}

func (e eventStoreEventService) Listen(subject v1.Subject) {
	data := v1.LogEvent{
		ProcessId: subject.Pipeline.ProcessId,
		Log:       subject.Log,
		Step:      subject.Step,
	}
	header := make(map[string]string)
	header["Content-Type"] = "application/json"
	header["token"] = config.Token
	b, err := json.Marshal(data)
	if err != nil {
		log.Println(err.Error())
		return
	}
	e.httpPublisher.Post(config.EventStoreUrl+"/logs", header, b)
}

// NewEventStoreLogEventService returns Observer type service
func NewEventStoreLogEventService(httpPublisher service.HttpClient) service.Observer {
	return &eventStoreEventService{
		httpPublisher: httpPublisher,
	}
}
