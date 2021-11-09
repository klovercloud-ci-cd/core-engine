package inmemory

import (
	"container/list"
	"encoding/json"
	v1 "github.com/klovercloud-ci/core/v1"
	"github.com/klovercloud-ci/core/v1/repository"
	"log"
)

type processEventRepository struct {
}

func (p processEventRepository) Store(data v1.ProcessEvent) {
	if ProcessEventStore == nil {
		ProcessEventStore = map[string]*list.List{}
	}
	_, ok := ProcessEventStore[data.ProcessId]
	if !ok {
		ProcessEventStore[data.ProcessId] = list.New()
	}
	ProcessEventStore[data.ProcessId].PushBack(data.Data)
}

func (p processEventRepository) GetByProcessId(processId string) map[string]interface{} {
	if _, ok := ProcessEventStore[processId]; ok {
		e := ProcessEventStore[processId]
		if ProcessEventStore[processId].Front() != nil {
			m := make(map[string]interface{})
			t := &e.Front().Value
			jsonString, err := json.Marshal(t)
			if err != nil {
				log.Println(err.Error())
			}
			err = json.Unmarshal(jsonString, &m)
			if err != nil {
				log.Println(err.Error())
			}
			return m

		}
	}
	return nil
}

func (p processEventRepository) DequeueByProcessId(processId string) map[string]interface{} {
	if _, ok := ProcessEventStore[processId]; ok {
		e := ProcessEventStore[processId]
		if ProcessEventStore[processId].Front() != nil {
			m := make(map[string]interface{})
			t := e.Remove(e.Front())
			jsonString, MarshalErr := json.Marshal(&t)
			if MarshalErr != nil {
				log.Println(MarshalErr.Error())
			}
			UnmarshalErr := json.Unmarshal(jsonString, &m)
			if UnmarshalErr != nil {
				log.Println(UnmarshalErr.Error())
			}
			return m
		}
	}
	return nil
}

// NewProcessEventRepository returns ProcessEventRepository type object
func NewProcessEventRepository() repository.ProcessEventRepository {
	return &processEventRepository{}
}
