package logic

import (
	v1 "github.com/klovercloud-ci/core/v1"
	"github.com/klovercloud-ci/core/v1/repository"
	"github.com/klovercloud-ci/core/v1/service"
)

type processEventService struct {
	repo repository.ProcessEventRepository
}

func (p processEventService) Listen(listener v1.Subject) {
	if listener.EventData!=nil{
		p.repo.Store(v1.PipelineProcessEvent{
			ProcessId: listener.Pipeline.ProcessId,
			Data:      listener.EventData,
		})
	}
}

func (p processEventService) Store(data v1.PipelineProcessEvent) {
	p.repo.Store(data)
}

func (p processEventService) GetByProcessId(processId string) map[string]interface{} {
	return p.repo.GetByProcessId(processId)
}

func (p processEventService) DequeueByProcessId(processId string) map[string]interface{} {
	return p.repo.DequeueByProcessId(processId)
}

func NewProcessEventService(repo repository.ProcessEventRepository) service.ProcessEvent {
	return &processEventService{
		repo: repo,
	}
}