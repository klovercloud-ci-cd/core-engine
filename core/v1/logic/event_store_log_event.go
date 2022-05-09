package logic

import (
	"encoding/json"
	"fmt"
	"github.com/klovercloud-ci-cd/core-engine/config"
	v1 "github.com/klovercloud-ci-cd/core-engine/core/v1"
	"github.com/klovercloud-ci-cd/core-engine/core/v1/service"
	"log"
	"strconv"
	"time"
)

type eventStoreEventService struct {
	httpPublisher service.HttpClient
}

func (e eventStoreEventService) Listen(subject v1.Subject) {
	if subject.Log == "" || subject.Step == "" {
		return
	}
	claim, _ := strconv.Atoi(fmt.Sprint(subject.EventData["claim"]))
	data := v1.LogEvent{
		ProcessId: subject.Pipeline.ProcessId,
		Log:       subject.Log,
		Step:      subject.Step,
		Footmark:  fmt.Sprint(subject.EventData["footmark"]),
		CreatedAt: time.Now().UTC(),
		Claim:     claim,
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
