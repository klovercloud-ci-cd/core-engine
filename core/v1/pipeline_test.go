package v1

import (
	"github.com/stretchr/testify/assert"
	"log"
	"reflect"
	"testing"
)

func Test_loadArgs(t *testing.T) {
	type TestCase struct {
		data Pipeline
		expected map[string]string
		actual map[string]string
	}
	variable:=Variable{}
	variable.ConfigMaps= append(variable.ConfigMaps, struct {
		Name string `json:"name"`
		Namespace   string `json:"namespace"`
	}{Name: "configMap0", Namespace:"klovercloud"})
	pipeline:=Pipeline{
		Steps: []Step{Step{
			Arg: variable,
		}},
	}
	testCase:=TestCase{
		data:    pipeline,
		expected: map[string]string{"env1":"value1","env2":"value2"},
	}
	pipeline.loadArgs(&MockK8sResource{})
	testCase.actual=pipeline.Steps[0].Arg.Data
	if !reflect.DeepEqual(testCase.expected, testCase.actual){
		assert.ElementsMatch(t, testCase.expected, testCase.actual)
	}
	variable.Secrets= append(variable.Secrets, struct {
		Name string `json:"name"`
		Namespace   string `json:"namespace"`
	}{Name: "secret0", Namespace:"klovercloud"})

	pipeline.Steps[0].Arg=variable
	testCase.expected=map[string]string{"env1":"value1","env2":"value2","key1":"value1","key2":"value2"}

	pipeline.loadArgs(&MockK8sResource{})
	testCase.actual=pipeline.Steps[0].Arg.Data
	if !reflect.DeepEqual(testCase.expected, testCase.actual){
		assert.ElementsMatch(t, testCase.expected, testCase.actual)
	}
}