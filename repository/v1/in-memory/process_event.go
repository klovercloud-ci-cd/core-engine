package in_memory

import (
	v1 "github.com/klovercloud-ci/core/v1"
	"github.com/klovercloud-ci/core/v1/repository"
	"github.com/orcaman/concurrent-map"
)

var processEvents cmap.ConcurrentMap


type processEventRepository struct {

}

func (p processEventRepository) Store(data v1.PipelineProcessStatus) {
	if processEvents==nil{
		processEvents=cmap.New()
	}
	processEvents.Set(data.ProcessId,data.Data)
}

func (p processEventRepository) GetByProcessId(processId string) map[string]interface{}{
	if tmp, ok := processEvents.Get(processId); ok {
		return tmp.(map[string]interface{})
	}
	return nil
}

func NewProcessEventRepository() (repository.ProcessEventRepository) {
	return &processEventRepository{
	}
}
