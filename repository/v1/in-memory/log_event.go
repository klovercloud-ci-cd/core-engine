package in_memory

import (
	v1 "github.com/klovercloud-ci/core/v1"
	"github.com/klovercloud-ci/core/v1/repository"
)

type logEventRepository struct {
}

func (l logEventRepository) Store(log v1.LogEvent)  {
	panic("implement me")
}

func (l logEventRepository) GetByBuildId(buildId string, option v1.LogEventQueryOption) []string {
	panic("implement me")
}

func NewLogEventRepository() (repository.LogEventRepository) {
	return &logEventRepository{
	}
}
