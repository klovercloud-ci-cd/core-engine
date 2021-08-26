package in_memory

import (
	v1 "github.com/klovercloud-ci/core/v1"
	"github.com/klovercloud-ci/core/v1/repository"
)

type logEventRepository struct {
}

func (l logEventRepository) Store(log v1.LogEvent)  {
	if IndexedLogEvents[log.ProcessId]==nil{
		IndexedLogEvents[log.ProcessId]=make(map[int32]v1.LogEvent)
	}
	IndexedLogEvents[log.ProcessId][log.Index]=log
}

func (l logEventRepository) GetByProcessId(processId string, option v1.LogEventQueryOption) []string {
	IndexedData:=IndexedLogEvents[processId]

	var data []string
	for i:=option.IndexFrom;i<=option.IndexTo;i++{
		data= append(data, IndexedData[i].Log)

	}
	return data
}

func NewLogEventRepository() (repository.LogEventRepository) {
	return &logEventRepository{
	}
}
