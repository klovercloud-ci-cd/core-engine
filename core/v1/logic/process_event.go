package logic

import (
	"fmt"
	v1 "github.com/klovercloud-ci-cd/core-engine/core/v1"
	"github.com/klovercloud-ci-cd/core-engine/core/v1/repository"
	"github.com/klovercloud-ci-cd/core-engine/core/v1/service"
)

type processEventService struct {
	repo repository.ProcessEventRepository
}

func (p processEventService) Listen(listener v1.Subject) {
	if listener.EventData != nil {
		p.repo.Store(v1.ProcessEvent{
			ProcessId: listener.Pipeline.ProcessId,
			Data:      listener.EventData,
			CompanyId: fmt.Sprint(listener.EventData["company_id"]),
		})
	}
}

func (p processEventService) Store(data v1.ProcessEvent) {
	p.repo.Store(data)
}

func (p processEventService) GetByCompanyId(companyId string) map[string]interface{} {
	return p.repo.GetByCompanyId(companyId)
}

func (p processEventService) DequeueByCompanyId(companyId string) map[string]interface{} {
	return p.repo.DequeueByCompanyId(companyId)
}

// NewProcessEventService returns ProcessEvent type service
func NewProcessEventService(repo repository.ProcessEventRepository) service.ProcessEvent {
	return &processEventService{
		repo: repo,
	}
}
