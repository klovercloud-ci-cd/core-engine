package logic

import (
	v1 "github.com/klovercloud-ci-cd/core-engine/core/v1"
	"github.com/klovercloud-ci-cd/core-engine/core/v1/repository"
	"github.com/klovercloud-ci-cd/core-engine/core/v1/service"
	"time"
)

type logEventService struct {
	repo repository.LogEventRepository
}

func (l logEventService) Listen(listener v1.Subject) {
	if listener.Log != "" {
		l.repo.Store(v1.LogEvent{
			ProcessId: listener.Pipeline.ProcessId,
			Log:       listener.Log,
			Step:      listener.Step,
			CreatedAt: time.Time{}.UTC(),
		})
	}
}

func (l logEventService) Store(log v1.LogEvent) {
	l.repo.Store(log)
}

func (l logEventService) GetByProcessId(processId string, option v1.LogEventQueryOption) ([]string, int64) {
	return l.repo.GetByProcessId(processId, option)
}

// NewLogEventService returns LogEvent type service
func NewLogEventService(repo repository.LogEventRepository) service.LogEvent {
	return &logEventService{
		repo: repo,
	}
}
