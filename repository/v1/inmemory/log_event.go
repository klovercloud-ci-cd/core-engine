package inmemory

import (
	v1 "github.com/klovercloud-ci/core/v1"
	"github.com/klovercloud-ci/core/v1/repository"
)

type logEventRepository struct {
}

func (l logEventRepository) Store(log v1.LogEvent) {
	if len(IndexedLogEvents) == 0 {
		IndexedLogEvents = make(map[string][]v1.LogEvent)
	}
	IndexedLogEvents[log.ProcessId] = append(IndexedLogEvents[log.ProcessId], log)
}

func (l logEventRepository) GetByProcessId(processId string, option v1.LogEventQueryOption) ([]string, int64) {
	logEvents := IndexedLogEvents[processId]
	var data []string
	for i := 0; i < len(logEvents); i++ {
		if option.Step != "" {
			if logEvents[i].Step == option.Step {
				data = append(data, logEvents[i].Log)
			}
		} else {
			data = append(data, logEvents[i].Log)
		}
	}
	logs := paginate(data, option.Pagination.Page, option.Pagination.Limit)
	return logs, int64(len(data))
}

func paginate(logs []string, page int64, limit int64) []string {
	if page < 0 || limit <= 0 {
		return nil
	}
	var startIndex, endIndex int64
	if page == 0 {
		startIndex = 0
	} else {
		startIndex = page * limit
	}
	endIndex = startIndex + limit
	if startIndex >= int64(len(logs)) {
		return nil
	}
	if endIndex >= int64(len(logs)) {
		return logs[startIndex:]
	}
	return logs[startIndex:endIndex]
}

// NewLogEventRepository returns LogEventRepository type object
func NewLogEventRepository() repository.LogEventRepository {
	return &logEventRepository{}
}
