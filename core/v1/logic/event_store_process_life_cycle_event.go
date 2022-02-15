package logic

import (
	"encoding/json"
	"fmt"
	"github.com/klovercloud-ci-cd/core-engine/api/common"
	"github.com/klovercloud-ci-cd/core-engine/config"
	"github.com/klovercloud-ci-cd/core-engine/core/v1"
	"github.com/klovercloud-ci-cd/core-engine/core/v1/service"
	"github.com/klovercloud-ci-cd/core-engine/enums"
	"log"
	"strconv"
	"strings"
	"time"
)

type eventStoreProcessLifeCycleService struct {
	httpPublisher service.HttpClient
}

func (e eventStoreProcessLifeCycleService) PullJenkinsJobStepsEvents() []v1.ProcessLifeCycleEvent {
	url := config.EventStoreUrl + "/process_life_cycle_events?count=" + strconv.FormatInt(config.AllowedConcurrentBuild, 10) + "&step_type=" + string(enums.JENKINS_JOB)
	header := make(map[string]string)
	header["Authorization"] = "token " + config.Token
	header["Accept"] = "application/json"
	header["token"] = config.Token
	data, err := e.httpPublisher.Get(url, header)
	if err != nil {
		// send to observer
		log.Println("[ERROR] Failed to connect event bank.", err.Error())
		return nil
	}
	response := common.ResponseDTO{}
	err = json.Unmarshal(data, &response)
	if err != nil {
		log.Println(err.Error())
		// send to observer
		return nil
	}
	b, err := json.Marshal(response.Data)
	if err != nil {
		log.Println(err.Error())
		// send to observer
		return nil
	}
	events := []v1.ProcessLifeCycleEvent{}
	err = json.Unmarshal(b, &events)
	if err != nil {
		log.Println(err.Error())
		// send to observer
		return nil
	}
	return events
}

func (e eventStoreProcessLifeCycleService) PullIntermediaryStepsEvents() []v1.ProcessLifeCycleEvent {
	url := config.EventStoreUrl + "/process_life_cycle_events?count=" + strconv.FormatInt(config.AllowedConcurrentBuild, 10) + "&step_type=" + string(enums.INTERMEDIARY)
	header := make(map[string]string)
	header["Authorization"] = "token " + config.Token
	header["Accept"] = "application/json"
	header["token"] = config.Token
	data, err := e.httpPublisher.Get(url, header)
	if err != nil {
		// send to observer
		log.Println("[ERROR] Failed to connect event bank.", err.Error())
		return nil
	}
	response := common.ResponseDTO{}
	err = json.Unmarshal(data, &response)
	if err != nil {
		log.Println(err.Error())
		// send to observer
		return nil
	}
	b, err := json.Marshal(response.Data)
	if err != nil {
		log.Println(err.Error())
		// send to observer
		return nil
	}
	events := []v1.ProcessLifeCycleEvent{}
	err = json.Unmarshal(b, &events)
	if err != nil {
		log.Println(err.Error())
		// send to observer
		return nil
	}
	return events
}

func (e eventStoreProcessLifeCycleService) PullBuildCancellingEvents() []v1.ProcessLifeCycleEvent {
	url := config.EventStoreUrl + "/process_life_cycle_events?count=step_type=" + string(enums.BUILD) + "&status=" + string(enums.REQUESTED_TO_CANCEL)
	header := make(map[string]string)
	header["Authorization"] = "token " + config.Token
	header["Accept"] = "application/json"
	header["token"] = config.Token
	data, err := e.httpPublisher.Get(url, header)
	if err != nil {
		// send to observer
		log.Println(err.Error())
		return nil
	}
	response := common.ResponseDTO{}
	err = json.Unmarshal(data, &response)
	if err != nil {
		log.Println(err.Error())
		// send to observer
		return nil
	}
	b, err := json.Marshal(response.Data)
	if err != nil {
		log.Println(err.Error())
		// send to observer
		return nil
	}
	events := []v1.ProcessLifeCycleEvent{}
	err = json.Unmarshal(b, &events)
	if err != nil {
		log.Println(err.Error())
		// send to observer
		return nil
	}
	return events
}

func (e eventStoreProcessLifeCycleService) PullBuildEvents() []v1.ProcessLifeCycleEvent {
	url := config.EventStoreUrl + "/process_life_cycle_events?count=" + strconv.FormatInt(config.AllowedConcurrentBuild, 10) + "&step_type=" + string(enums.BUILD)
	header := make(map[string]string)
	header["Authorization"] = "token " + config.Token
	header["Accept"] = "application/json"
	header["token"] = config.Token
	data, err := e.httpPublisher.Get(url, header)
	if err != nil {
		// send to observer
		log.Println("[ERROR] Failed to connect event bank.", err.Error())
		return nil
	}
	response := common.ResponseDTO{}
	err = json.Unmarshal(data, &response)
	if err != nil {
		log.Println(err.Error())
		// send to observer
		return nil
	}
	b, err := json.Marshal(response.Data)
	if err != nil {
		log.Println(err.Error())
		// send to observer
		return nil
	}
	events := []v1.ProcessLifeCycleEvent{}
	err = json.Unmarshal(b, &events)
	if err != nil {
		log.Println(err.Error())
		// send to observer
		return nil
	}
	return events
}

func (e eventStoreProcessLifeCycleService) Listen(subject v1.Subject) {
	if subject.EventData["status"] == nil {
		return
	}
	data := []v1.ProcessLifeCycleEvent{}
	var processLifeCycleEvent v1.ProcessLifeCycleEvent
	if subject.EventData["status"] != enums.NON_INITIALIZED {
		var nextSteps []string
		for _, each := range subject.Pipeline.Steps {
			if each.Name == subject.Step {
				processLifeCycleEvent = v1.ProcessLifeCycleEvent{
					ProcessId: subject.Pipeline.ProcessId,
					Step:      subject.Step,
					StepType:  enums.STEP_TYPE(fmt.Sprintf("%v", subject.EventData["type"])),
					Next:      strings.Split(fmt.Sprintf("%v", subject.EventData["next"]), ","),
				}
				nextSteps = each.Next

			}
		}
		if subject.EventData["status"] == string(enums.STEP_FAILED) || subject.EventData["status"] == string(enums.ERROR) || subject.EventData["status"] == string(enums.TERMINATING) {
			processLifeCycleEvent.Status = enums.FAILED
			data = append(data, processLifeCycleEvent)
		} else if subject.EventData["status"] == string(enums.SUCCESSFUL) {
			processLifeCycleEvent.Status = enums.COMPLETED
			data = append(data, processLifeCycleEvent)
			for _, each := range nextSteps {
				data = append(data, v1.ProcessLifeCycleEvent{
					ProcessId: subject.Pipeline.ProcessId,
					Status:    enums.PAUSED,
					Step:      each,
				})
			}

		}

	} else {

		processLifeCycleEvent := v1.ProcessLifeCycleEvent{
			ProcessId: subject.Pipeline.ProcessId,
			Status:    enums.NON_INITIALIZED,
			Step:      subject.Step,
			StepType:  enums.STEP_TYPE(fmt.Sprintf("%v", subject.EventData["type"])),
			Next:      strings.Split(fmt.Sprintf("%v", subject.EventData["next"]), ","),
			Agent:     fmt.Sprintf("%v", subject.EventData[string(enums.AGENT)]),
			Pipeline:  &subject.Pipeline,
			CreatedAt: time.Now().UTC(),
			Trigger:   enums.TRIGGER(fmt.Sprintf("%v", subject.EventData["trigger"])),
		}
		data = append(data, processLifeCycleEvent)

		for i, each := range subject.Pipeline.Steps {
			if i == 0 {
				continue
			}
			data = append(data, v1.ProcessLifeCycleEvent{
				ProcessId: subject.Pipeline.ProcessId,
				Step:      each.Name,
				Status:    enums.NON_INITIALIZED,
				StepType:  each.Type,
				Next:      each.Next,
				Agent:     each.Params[enums.AGENT],
				Pipeline:  nil,
				CreatedAt: time.Now().UTC(),
				Trigger:   each.Trigger,
			})
		}
	}
	type ProcessLifeCycleEventList struct {
		Events []v1.ProcessLifeCycleEvent `bson:"events" json :"events"`
	}
	if len(data) > 0 {
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

// NewEventStoreProcessLifeCycleService returns ProcessLifeCycleEvent type service
func NewEventStoreProcessLifeCycleService(httpPublisher service.HttpClient) service.ProcessLifeCycleEvent {
	return &eventStoreProcessLifeCycleService{
		httpPublisher: httpPublisher,
	}
}
