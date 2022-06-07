package logic

import (
	"fmt"
	"github.com/klovercloud-ci-cd/core-engine/config"
	v1 "github.com/klovercloud-ci-cd/core-engine/core/v1"
	"github.com/klovercloud-ci-cd/core-engine/core/v1/repository"
	"github.com/klovercloud-ci-cd/core-engine/core/v1/service"
	"github.com/klovercloud-ci-cd/core-engine/enums"
	"strings"
	"time"
)

type processLifeCycleEventService struct {
	repo repository.ProcessLifeCycleEventRepository
}

func (p processLifeCycleEventService) PullBuildCancellingEvents() []v1.ProcessLifeCycleEvent {
	return p.repo.PullCancellingStepsByProcessStatusAndStepType(string(enums.BUILD))
}

func (p processLifeCycleEventService) PullJenkinsJobStepsEvents() []v1.ProcessLifeCycleEvent {
	return p.PullNonInitializedAndAutoTriggerEnabledEventsByStepType(config.AllowedConcurrentBuild, string(enums.JENKINS_JOB))
}

func (p processLifeCycleEventService) PullIntermediaryStepsEvents() []v1.ProcessLifeCycleEvent {
	return p.PullNonInitializedAndAutoTriggerEnabledEventsByStepType(config.AllowedConcurrentBuild, string(enums.INTERMEDIARY))
}

func (p processLifeCycleEventService) PullBuildEvents() []v1.ProcessLifeCycleEvent {
	return p.PullNonInitializedAndAutoTriggerEnabledEventsByStepType(config.AllowedConcurrentBuild, string(enums.BUILD))
}

func (p processLifeCycleEventService) Listen(subject v1.Subject) {
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
		if subject.EventData["status"] == string(enums.FAILED) || subject.EventData["status"] == string(enums.STEP_FAILED) || subject.EventData["status"] == string(enums.ERROR) || subject.EventData["status"] == string(enums.TERMINATING) {
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
	p.Store(data)
}

func (p processLifeCycleEventService) PullNonInitializedAndAutoTriggerEnabledEventsByStepType(count int64, stepType string) []v1.ProcessLifeCycleEvent {
	return p.repo.PullNonInitializedAndAutoTriggerEnabledEventsByStepType(count, stepType)
}

func (p processLifeCycleEventService) PullPausedAndAutoTriggerEnabledResourcesByAgentName(count int64, agent string) []v1.AgentDeployableResource {
	resources := []v1.AgentDeployableResource{}
	events := p.repo.PullPausedAndAutoTriggerEnabledResourcesByAgentName(count, agent)
	for _, event := range events {

		var step *v1.Step
		for _, each := range event.Pipeline.Steps {
			if each.Name == event.Step {
				step = &each
				break
			}
		}
		if step != nil {
			resources = append(resources, v1.AgentDeployableResource{
				Step:        step.Name,
				ProcessId:   event.ProcessId,
				Descriptors: step.Descriptors,
				Type:        enums.PIPELINE_RESOURCE_TYPE(step.Params["type"]),
				Name:        step.Params["name"],
				Namespace:   step.Params["namespace"],
				Images:      strings.Split(fmt.Sprintf("%v", step.Params["images"]), ","),
			})
		}
	}
	return resources
}

func (p processLifeCycleEventService) Store(events []v1.ProcessLifeCycleEvent) {
	p.repo.Store(events)
}

// NewProcessLifeCycleEventService returns ProcessLifeCycleEvent type service
func NewProcessLifeCycleEventService(repo repository.ProcessLifeCycleEventRepository) service.ProcessLifeCycleEvent {
	return &processLifeCycleEventService{
		repo: repo,
	}
}
