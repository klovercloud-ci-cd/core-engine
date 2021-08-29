package in_memory

import (
	v1 "github.com/klovercloud-ci/core/v1"
	"github.com/klovercloud-ci/core/v1/repository"
)

type logEventRepository struct {
}

func (l logEventRepository) Store(log v1.LogEvent)  {
	IndexedLogEvents[log.ProcessId]= append(IndexedLogEvents[log.ProcessId], log)
}

func (l logEventRepository) GetByProcessId(processId string, option v1.LogEventQueryOption) []string {
	logEvents:=IndexedLogEvents[processId]

	var data []string
	for i:=0;i<len(logEvents);i++{
		data= append(data, logEvents[i].Log)

	}
	return data
}

func NewLogEventRepository() (repository.LogEventRepository) {
	return &logEventRepository{
	}
}
