package mongo

import (
	v1 "github.com/klovercloud-ci/core/v1"
	"time"
)

var data []v1.LogEvent

func InitData() {
	data = []v1.LogEvent{
		{
			Index:     1,
			BuildId:   "6652",
			Log:       "hi",
			Step:      "first",
			CreatedAt: time.Now(),

		},
		{
			Index:     2,
			BuildId:   "5231",
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
