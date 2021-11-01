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
		data     v1.LogEvent
		expected int
		actual   int
		option   v1.LogEventQueryOption
	}
	var testCases []TestData

	for i := 0; i < 50; i++ {
		testCases = append(testCases, TestData{
			data: v1.LogEvent{
				ProcessId: "01",
				Log:       "log-" + strconv.Itoa(i),
				Step:      "BUILD",
				CreatedAt: time.Time{}.UTC(),
			},
			expected: i + 1,
			option: v1.LogEventQueryOption{
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
	logEventService := NewLogEventService(in_memory.NewLogEventRepository())
	for _, each := range testCases {
		logEventService.Store(each.data)
		logs, _ := logEventService.GetByProcessId(each.data.ProcessId, each.option)
		if !reflect.DeepEqual(len(logs), each.expected) {
			assert.ElementsMatch(t, len(logs), each.expected)
		}
	}
	testCases = append(testCases, TestData{
		data: v1.LogEvent{
			ProcessId: "01",
			Log:       "log-99",
			Step:      "DEVELOP",
			CreatedAt: time.Time{}.UTC(),
		},
		expected: 1,
		option: v1.LogEventQueryOption{
			Pagination: struct {
				Page  int64
				Limit int64
			}{
				Page:  0,
				Limit: 10,
			},
			Step: "DEVELOP",
		},
	})
	logEventService = NewLogEventService(in_memory.NewLogEventRepository())
	logEventService.Store(testCases[50].data)
	logs, _ := logEventService.GetByProcessId(testCases[50].data.ProcessId, testCases[50].option)
	if !reflect.DeepEqual(len(logs), testCases[50].expected) {
		assert.ElementsMatch(t, len(logs), testCases[50].expected)
	}
}

func Test_Store(t *testing.T) {
	type TestData struct {
		data     v1.LogEvent
		expected int64
		actual   int
		option   v1.LogEventQueryOption
	}
	var testCases []TestData

	for i := 0; i < 50; i++ {
		testCases = append(testCases, TestData{
			data: v1.LogEvent{
				ProcessId: "01",
				Log:       "log-" + strconv.Itoa(i),
				Step:      "BUILD",
				CreatedAt: time.Time{}.UTC(),
			},
			expected: int64(i + 52),
			option: v1.LogEventQueryOption{
				Pagination: struct {
					Page  int64
					Limit int64
				}{
					Page:  0,
					Limit: 100,
				},
			},
		})
	}
	logEventService := NewLogEventService(in_memory.NewLogEventRepository())
	for _, each := range testCases {
		logEventService.Store(each.data)
		_, size := logEventService.GetByProcessId(each.data.ProcessId, each.option)
		if size != each.expected {
			assert.ElementsMatch(t, size, each.expected)
		}
	}
}
func Test_Listen(t *testing.T) {
	type TestData struct {
		data     v1.Subject
		expected int64
	}
	var testCases []TestData
	count := int64(0)
	for i := 0; i < 50; i++ {
		log := "testlog"
		if i%2 == 0 {
			log = ""
		}
		if i%2 != 0 {
			count = count + 1
		}
		testCases = append(testCases, TestData{
			data: v1.Subject{
				Pipeline: v1.Pipeline{
					ApiVersion: "123",
					Name:       "test1",
					ProcessId:  "1",
				},
				Log:  log,
				Step: "Build",
			},
			expected: count,
		})
	}
	logEventService := NewLogEventService(in_memory.NewLogEventRepository())
	for _, each := range testCases {
		logEventService.Listen(each.data)
		_, size := logEventService.GetByProcessId(each.data.Pipeline.ProcessId, v1.LogEventQueryOption{
			Pagination: struct {
				Page  int64
				Limit int64
			}{0, 100},
		})
		if size != each.expected {
			assert.ElementsMatch(t, size, each.expected)
		}
	}
}
