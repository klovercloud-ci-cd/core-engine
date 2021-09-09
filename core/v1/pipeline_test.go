package v1

import (
	"github.com/klovercloud-ci/enums"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestPipeline_Validate(t *testing.T) {
	type TestCase struct {
		data Pipeline
		expected string
		actual string
	}
	pipeline := Pipeline{
		ApiVersion: "default",
		Name:       "test",
		ProcessId:  "0125632",
		Label:      map[string]string{"env1": "value1", "env2": "value2"},
		Steps:      []Step{{
			Name: "test",
			Type: enums.BUILD,
			ServiceAccount: "Testacc",
			Input:          Resource{
				Type: "git",
				Url: "github.com/abc",
				Revision: "12356",
			},
			Outputs: []Resource{{
				Type: "git",
				Url: "github.com/abc",
				Revision: "12356",
			}},
		}},
	}
	testcases:=[]TestCase{}
	testcases= append(testcases, TestCase{
		data:     pipeline,
		expected: "",
	})
	pipeline.ApiVersion = ""
	testcases =  append(testcases, TestCase{
		data:     pipeline,
		expected: "Api version is required!",
	})
	pipeline.ApiVersion = "default"
	pipeline.Name = ""
	testcases = append(testcases, TestCase{
		data:     pipeline,
		expected: "Pipeline name is required!",
	})
	pipeline.Name = "test"
	pipeline.ProcessId = ""
	testcases = append(testcases, TestCase{
		data:     pipeline,
		expected: "Pipeline process id is required!",
	})
	for i,_:=range testcases{
		err:=Pipeline.Validate(testcases[i].data)
		if err!=nil{
			testcases[i].actual =err.Error()
		}else{
			testcases[i].actual=""
		}
		if !reflect.DeepEqual(testcases[i].expected, testcases[i].actual) {
			assert.ElementsMatch(t, testcases[i].expected, testcases[i].actual)
		}
	}
}