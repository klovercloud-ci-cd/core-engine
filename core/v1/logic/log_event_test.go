package logic

import (
	v1 "github.com/klovercloud-ci/core/v1"
	in_memory "github.com/klovercloud-ci/repository/v1/in-memory"
	"github.com/stretchr/testify/assert"
	"reflect"
	"strconv"
	"testing"
	"time"
)

func Test_GetByProcessId(t *testing.T) {
	type TestData struct {
		data                  v1.LogEvent
		expected int
		actual int
		option v1.LogEventQueryOption
	}
	var testCases [] TestData

	for i:=0;i<50;i++{
		testCases= append(testCases, TestData{
			data: v1.LogEvent{
				ProcessId: "01",
				Log:       "log-" + strconv.Itoa(i),
				Step:      "BUILD",
				CreatedAt: time.Time{}.UTC(),
			},
			expected: i + 1,
			option:   v1.LogEventQueryOption{
				Pagination: struct {
					Page  int64
					Limit int64
				}{
					Page:  0,
					Limit: int64(i + 1),
				},
			},
		})
	}
	logEventService:=NewLogEventService(in_memory.NewLogEventRepository())
	for _,each:=range testCases{
		logEventService.Store(each.data)
		logs,_:=logEventService.GetByProcessId(each.data.ProcessId,each.option)
		if !reflect.DeepEqual(len(logs), each.expected){
			assert.ElementsMatch(t, len(logs), each.expected)
		}
	}


}