package in_memory

import (
	v1 "github.com/klovercloud-ci/core/v1"
	"github.com/klovercloud-ci/core/v1/repository"
)

type logEventRepository struct {
}

func (l logEventRepository) Store(log v1.LogEvent)  {

	if IndexedLogEvents[log.BuildId]==nil{
		IndexedLogEvents[log.BuildId]=make(map[int32]v1.LogEvent)
	}
	IndexedLogEvents[log.BuildId][log.Index]=log
}

func (l logEventRepository) GetByBuildId(buildId string, option v1.LogEventQueryOption) []string {
	IndexedData:=IndexedLogEvents[buildId]

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
