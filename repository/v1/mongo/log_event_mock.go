package mongo

import (
	"context"
	v1 "github.com/klovercloud-ci/core/v1"
	"github.com/klovercloud-ci/core/v1/repository"
	"time"
)

var data []v1.LogEvent

func InitData() {
	data = []v1.LogEvent{
		{
			Index:     1,
			ProcessId: "6652",
			Log:       "hi",
			Step:      "first",
			CreatedAt: time.Now(),

		},
		{
			Index:     2,
			ProcessId: "5231",
			Log:       "hello",
			Step:      "second",
			CreatedAt: time.Now(),

		},
	}
	this := new(logEventRepository)
	for _,each:=range data{
		this.Store(each)
	}
}
func NewMockLogEventRepository() repository.LogEventRepository{
	manager:=GetMockDmManager()
	manager.Db.Drop(context.Background())
	return &logEventRepository{
		manager: GetMockDmManager(),
		timeout: 3000,
	}

}