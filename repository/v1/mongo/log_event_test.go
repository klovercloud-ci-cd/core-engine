package mongo

import (
	"github.com/joho/godotenv"
	v1 "github.com/klovercloud-ci/core/v1"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"path"
	"reflect"
	"testing"
	"time"
)

func loadEnv(t *testing.T) error {
	dirname, err := os.Getwd()
	if err != nil {
		log.Println("ERROR:", err.Error())
		t.Fail()
	}
	dir, err := os.Open(path.Join(dirname, "../../../"))
	if err != nil {
		log.Println("ERROR:", err.Error())
		t.Fail()
	}
	err = godotenv.Load(os.ExpandEnv(dir.Name() + "/.env.mongo.test"))
	if err != nil {
		log.Println("ERROR:", err.Error())
		t.Fail()
	}
	return err
}

func TestLogEventRepository_Store(t *testing.T) {
	err:=loadEnv(t)
	if err!=nil{
		log.Println(err.Error())
	}
	type TestCase struct {
		data v1.LogEvent
		expected []string
		actual []string
	}
	testCases:=[]TestCase{}
	testCases= append(testCases, TestCase{
		data:     v1.LogEvent{
			ProcessId: "1",
			Log:       "Initializing pod",
			Step:      "buildImage",
			CreatedAt: time.Time{},
		},
		expected: []string{"Initializing pod"},
		actual: nil,
	})
	testCases= append(testCases, TestCase{
		data:     v1.LogEvent{
			ProcessId: "1",
			Log:       "Pulling Image",
			Step:      "buildImage",
			CreatedAt: time.Time{},
		},
		expected: []string{"Initializing pod","Pulling Image"},
		actual: nil,
	})

	l:=NewMockLogEventRepository()
	for _,each:=range testCases{
		l.Store(each.data)
		each.actual,_=l.GetByProcessId(each.data.ProcessId,v1.LogEventQueryOption{})
		if !reflect.DeepEqual(each.expected, each.actual){
			assert.ElementsMatch(t, each.expected, each.actual)
		}
	}
}

func TestLogEventRepository_GetByProcessId(t *testing.T) {
	type TestCase struct {
		processId string
		option v1.LogEventQueryOption
		expected []string
		actual []string
	}
	testCases:=[]TestCase{}
	testCases= append(testCases, TestCase{
		processId: "1",
		option:    v1.LogEventQueryOption{
			Pagination: struct {
				Page  int64
				Limit int64
			}{
				Page:  0,
				Limit: 0,
			},
		},
		expected:  []string{"Initializing pod","Pulling Image"},
	})
	testCases= append(testCases, TestCase{
		processId: "1",
		option:    v1.LogEventQueryOption{
			Pagination: struct {
				Page  int64
				Limit int64
			}{
				Page:  0,
				Limit: 1,
			},
		},
		expected:  []string{"Initializing pod"},
	})
	testCases= append(testCases, TestCase{
		processId: "1",
		option:    v1.LogEventQueryOption{
			Pagination: struct {
				Page  int64
				Limit int64
			}{
				Page:  1,
				Limit: 1,
			},
		},
		expected:  []string{"Pulling Image"},
	})
	testCases= append(testCases, TestCase{
		processId: "1",
		option:    v1.LogEventQueryOption{
			Pagination: struct {
				Page  int64
				Limit int64
			}{
				Page:  0,
				Limit: 2,
			},
		},
		expected:  []string{"Initializing pod","Pulling Image"},
	})
	testCases= append(testCases, TestCase{
		processId: "1",
		option:    v1.LogEventQueryOption{
			Pagination: struct {
				Page  int64
				Limit int64
			}{
				Page:  0,
				Limit: 3,
			},
		},
		expected:  []string{"Initializing pod","Pulling Image"},
	})

	testCases= append(testCases, TestCase{
		processId: "2",
		option:    v1.LogEventQueryOption{
			Step: "buildImage",
			Pagination: struct {
				Page  int64
				Limit int64
			}{Page:0 , Limit:10 },
		},
		expected:  []string{"Failed to initialize pod"},
	})

	testCases= append(testCases, TestCase{
		processId: "2",
		option:    v1.LogEventQueryOption{
			Step: "deployImage",
			Pagination: struct {
				Page  int64
				Limit int64
			}{Page:0 , Limit:1 },
		},
		expected:  []string{"Initializing pod"},
	})
	l:=NewMockLogEventRepository()
	data:= InitLogEventData()
	for _,each:=range data{
		l.Store(each)
	}
	for _,each:=range testCases{
		each.actual,_=l.GetByProcessId(each.processId,each.option)
		if !reflect.DeepEqual(each.expected, each.actual){
			assert.ElementsMatch(t, each.expected, each.actual)
		}
	}
}