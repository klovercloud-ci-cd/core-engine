package logic

import (
	v1 "github.com/klovercloud-ci/core/v1"
	"github.com/klovercloud-ci/core/v1/repository"
	"github.com/klovercloud-ci/core/v1/service"
)

type logEventService struct {
	repo repository.LogEventRepository
}

func (l logEventService) Store(log v1.LogEvent) {
	l.repo.Store(log)
}

func (l logEventService) GetByProcessId(processId string, option v1.LogEventQueryOption) ([]string,int64) {
	return l.repo.GetByProcessId(processId,option)
}

func NewLogEventService(repo repository.LogEventRepository) service.LogEvent {
	return &logEventService{
		repo: repo,
	}
}