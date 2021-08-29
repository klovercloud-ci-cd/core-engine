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
		each.actual=l.GetByProcessId(each.data.ProcessId,v1.LogEventQueryOption{})
		if !reflect.DeepEqual(each.expected, each.actual){
			assert.ElementsMatch(t, each.expected, each.actual)
		}
	}
}
