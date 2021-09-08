package in_memory

import (
	v1 "github.com/klovercloud-ci/core/v1"
	"github.com/stretchr/testify/assert"
	"log"
	"strconv"
	"testing"
	"time"
)

func initLogEvents(){
	IndexedLogEvents=make(map[string][]v1.LogEvent)
	IndexedLogEvents["01"]=[]v1.LogEvent{}
	IndexedLogEvents["02"]=[]v1.LogEvent{}
	for i := 0; i < 50; i++ {
		step:="BUILD"
		if i%2==0{
			step="DEPLOY"
		}
		IndexedLogEvents["01"] = append(IndexedLogEvents["01"],v1.LogEvent{
			ProcessId: "01",
			Log:       "log event "+strconv.Itoa(i),
			Step:      step,
			CreatedAt: time.Time{}.UTC(),
		})
		IndexedLogEvents["02"] = append(IndexedLogEvents["02"],v1.LogEvent{
			ProcessId: "02",
			Log:       "log event "+strconv.Itoa(i),
			Step:      step,
			CreatedAt: time.Time{}.UTC(),
		})
	}
}

func Test_GetByProcessId(t *testing.T) {
	initLogEvents()
	type TestData struct {
		Data   map[string][]v1.LogEvent
		Page, Limit, Expected int64
		ProcessId,step string
	}

	data := []TestData{
		{Data: IndexedLogEvents, Page: 0, Limit: 6, Expected: 6,ProcessId: "01"},
		{Data: IndexedLogEvents, Page: 1, Limit: 6, Expected: 6,ProcessId: "01"},
		{Data: IndexedLogEvents, Page: 2, Limit: 6, Expected: 6,ProcessId: "01"},
		{Data: IndexedLogEvents, Page: 3, Limit: 6, Expected: 6,ProcessId: "01"},
		{Data: IndexedLogEvents, Page: 4, Limit: 6, Expected: 6,ProcessId: "01"},
		{Data: IndexedLogEvents, Page: 5, Limit: 6, Expected: 6,ProcessId: "01"},
		{Data: IndexedLogEvents, Page: 6, Limit: 6, Expected: 6,ProcessId: "01"},
		{Data: IndexedLogEvents, Page: 7, Limit: 6, Expected: 6,ProcessId: "01"},
		{Data: IndexedLogEvents, Page: 8, Limit: 6, Expected: 2,ProcessId: "01"},
		{Data: IndexedLogEvents, Page: 9, Limit: 6, Expected: 0,ProcessId: "01"},
		{Data: IndexedLogEvents, Page: 0, Limit: 50, Expected: 50,ProcessId: "01"},
		{Data: IndexedLogEvents, Page: 0, Limit: 51, Expected: 50,ProcessId: "01"},
		{Data: IndexedLogEvents, Page: 1, Limit: 25, Expected: 25,ProcessId: "01"},
		{Data: IndexedLogEvents,Page: -1, Limit: 25, Expected: 0,ProcessId: "01"},
		{Data: IndexedLogEvents,Page: 1, Limit: 0, Expected: 0,ProcessId: "01"},
		{Data: IndexedLogEvents, Page: 0, Limit: 100, Expected: 25,ProcessId: "01",step: "BUILD"},
	}

	repo:=NewLogEventRepository()
	for _, each := range data {
		data,_:=repo.GetByProcessId(each.ProcessId,v1.LogEventQueryOption{
			Pagination: struct {
				Page  int64
				Limit int64
			}{
				Page:  each.Page,
				Limit: each.Limit,
			},
			Step: each.step,
		})
		if len(data) !=int(each.Expected) {
			assert.ElementsMatch(t,  len(data), int(each.Expected))
		}
	}

}

func Test_Store(t *testing.T) {
	initLogEvents()
	type TestData struct {
		Data   v1.LogEvent
		Expected int64
		ProcessId string
	}
	data := []TestData{
		{Data: v1.LogEvent{"01","test","BUILD",time.Time{}.UTC()},Expected: 51, ProcessId: "01"},
		{Data: v1.LogEvent{"01","test","BUILD",time.Time{}.UTC()},Expected: 52, ProcessId: "01"},
		{Data: v1.LogEvent{"02","test","BUILD",time.Time{}.UTC()},Expected: 51, ProcessId: "02"},
	}
	repo:=NewLogEventRepository()
	for i, each := range data {
		repo.Store(each.Data)
		_,size:=repo.GetByProcessId(data[i].ProcessId,v1.LogEventQueryOption{
			Pagination: struct {
				Page  int64
				Limit int64
			}{
				Page:  0,
				Limit: 100,
			},
		})
		log.Println(size, " ",each.Expected,size != each.Expected )
		if size != each.Expected {
			assert.ElementsMatch(t,  size, each.Expected)
		}
	}

}
func initDataForPagination() []string {
	var data []string
	for i := 0; i < 50; i++ {
		data = append(data, "log event "+strconv.Itoa(i))
	}
	return data
}
func Test_pagination(t *testing.T) {
	records := initDataForPagination()
	type TestData struct {
		Data                  []string
		Page, Limit, Expected int64
	}

	data := []TestData{
		{Data: records, Page: 0, Limit: 6, Expected: 6},
		{Data: records, Page: 1, Limit: 6, Expected: 6},
		{Data: records, Page: 2, Limit: 6, Expected: 6},
		{Data: records, Page: 3, Limit: 6, Expected: 6},
		{Data: records, Page: 4, Limit: 6, Expected: 6},
		{Data: records, Page: 5, Limit: 6, Expected: 6},
		{Data: records, Page: 6, Limit: 6, Expected: 6},
		{Data: records, Page: 7, Limit: 6, Expected: 6},
		{Data: records, Page: 8, Limit: 6, Expected: 2},
		{Data: records, Page: 9, Limit: 6, Expected: 0},
		{Data: records, Page: 0, Limit: 50, Expected: 50},
		{Data: records, Page: 0, Limit: 51, Expected: 50},
		{Data: records, Page: 1, Limit: 25, Expected: 25},
		{Data: records,Page: -1, Limit: 25, Expected: 0},
		{Data: records,Page: 1, Limit: 0, Expected: 0},
	}

	for _, each := range data {

		if len(paginate(each.Data, each.Page, each.Limit)) != int(each.Expected) {
			assert.ElementsMatch(t,  len(paginate(each.Data, each.Page, each.Limit)), each.Expected)
		}
	}

}