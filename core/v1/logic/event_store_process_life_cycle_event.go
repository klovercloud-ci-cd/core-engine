package logic

import (
	"encoding/json"
	"fmt"
	"github.com/klovercloud-ci/config"
	"github.com/klovercloud-ci/core/v1"
	"github.com/klovercloud-ci/core/v1/service"
	"github.com/klovercloud-ci/enums"
	"log"
	"strings"
	"time"
)

type eventStoreProcessLifeCycleService struct {
	httpPublisher service.HttpPublisher
}

func (e eventStoreProcessLifeCycleService) Listen(subject v1.Subject) {
	if subject.EventData["status"]==nil{
		return
	}
	if subject.EventData != nil {
		data := []v1.ProcessLifeCycleEvent{}
		processLifeCycleEvent := v1.ProcessLifeCycleEvent{
			ProcessId: subject.Pipeline.ProcessId,
			Step:      subject.Step,
			StepType:  enums.STEP_TYPE(fmt.Sprintf("%v", subject.EventData["type"])),
			Next:      strings.Split(fmt.Sprintf("%v", subject.EventData["next"]), ","),
			Agent:     fmt.Sprintf("%v", subject.EventData[string(enums.AGENT)]),
			Pipeline:  &subject.Pipeline,
			CreatedAt: time.Now().UTC(),
			Trigger:   enums.TRIGGER(fmt.Sprintf("%v", subject.EventData["trigger"])),
		}
		nestStepNameMap := make(map[string]bool)
		for _, name := range processLifeCycleEvent.Next {
			nestStepNameMap[name] = true
		}
		if subject.EventData["status"] == enums.INITIALIZING {
			processLifeCycleEvent.Status = enums.ACTIVE
			data = append(data,processLifeCycleEvent)
			for _, each := range subject.Pipeline.Steps {
				if _, ok := nestStepNameMap[each.Name]; ok {
					data = append(data, v1.ProcessLifeCycleEvent{
						ProcessId: subject.Pipeline.ProcessId,
						Step:      each.Name,
						Status:    enums.PAUSED,
						StepType: each.Type,
						Next:      each.Next,
						Agent:     each.Params[enums.AGENT],
						Pipeline:  nil,
						CreatedAt: time.Now().UTC(),
						Trigger:   each.Trigger,
					})
				}
			}

		}else if subject.EventData["status"] == string(enums.BUILD_FAILED) || subject.EventData["status"] == string(enums.ERROR) || subject.EventData["status"] == string(enums.TERMINATING){
			processLifeCycleEvent.Status = enums.FAILED
			data = append(data,processLifeCycleEvent)
		}else if subject.EventData["status"] == string(enums.SUCCESSFUL){
			processLifeCycleEvent.Status = enums.COMPLETED
			data = append(data,processLifeCycleEvent)
		}
		type ProcessLifeCycleEventList struct {
			Events [] v1.ProcessLifeCycleEvent `bson:"events" json :"events"`
		}
		if len(data)>0 {
			events := ProcessLifeCycleEventList{data}
			header := make(map[string]string)
			header["Content-Type"] = "application/json"
			b, err := json.Marshal(events)
			if err != nil {
				log.Println(err.Error())
				return
			}
			e.httpPublisher.Post(config.EventStoreUrl+"/process_life_cycle_events", header, b)
		}
	}
}

func NewEventStoreProcessLifeCycleService(httpPublisher service.HttpPublisher) service.Observer {
	return &eventStoreProcessLifeCycleService{
		httpPublisher: httpPublisher,
	}
}
