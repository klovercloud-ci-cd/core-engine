package in_memory

import (
	"container/list"
	v1 "github.com/klovercloud-ci/core/v1"
	"github.com/klovercloud-ci/core/v1/repository"
)


var processEventStore map[string]*list.List

type processEventRepository struct {

}



func (p processEventRepository) Store(data v1.PipelineProcessEvent) {
	if processEventStore==nil{
		processEventStore= map[string]*list.List{}
	}
	_,ok:=processEventStore[data.ProcessId]
	if !ok{
		processEventStore[data.ProcessId]=list.New()
	}
	processEventStore[data.ProcessId].PushBack(&data.Data)
}

func (p processEventRepository) GetByProcessId(processId string) map[string]interface{}{
	if _, ok := processEventStore[processId]; ok {
		e:=processEventStore[processId]
		if processEventStore[processId].Front()!=nil {
			t:=e.Front().Value
			return  t.(map[string]interface{})
		}
	}

	return nil
}

func (p processEventRepository) DequeueByProcessId(processId string) map[string]interface{} {
	if _, ok := processEventStore[processId]; ok {
		e:=processEventStore[processId]
		if processEventStore[processId].Front()!=nil {
			t:=e.Remove(e.Front())
			return  t.(map[string]interface{})
		}
	}

	return nil
}
func NewProcessEventRepository() (repository.ProcessEventRepository) {
	return &processEventRepository{
	}
}
