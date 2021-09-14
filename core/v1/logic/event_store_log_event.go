package logic

import (
	"encoding/json"
	"github.com/klovercloud-ci/config"
	v1 "github.com/klovercloud-ci/core/v1"
	"github.com/klovercloud-ci/core/v1/service"
	"log"
)

type eventStoreEventService struct {
	httpPublisher service.HttpPublisher
}

func (e eventStoreEventService) Listen(subject v1.Subject) {
	data:=v1.LogEvent{
		ProcessId: subject.Pipeline.ProcessId,
		Log:       subject.Log,
		Step:      subject.Step,
	}
	header:=make(map[string]string)
	header["Content-Type"]="application/json"
	b, err := json.Marshal(data)
	if err!=nil{
		log.Println(err.Error())
		return
	}
	e.httpPublisher.Post(config.EventStoreUrl+"/logs",header,b)
}

func NewEventStoreLogEventService(httpPublisher service.HttpPublisher) service.EventStoreLogEvent {
	return &eventStoreEventService{
		httpPublisher: httpPublisher,
	}
}