package logic

import (
	v1 "github.com/klovercloud-ci/core/v1"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func Test_loadArgs(t *testing.T) {
	type TestCase struct {
		data     v1.Pipeline
		expected map[string]string
		actual   map[string]string
	}
	variable := v1.Variable{}
	variable.ConfigMaps = append(variable.ConfigMaps, struct {
		Name      string `json:"name"`
		Namespace string `json:"namespace"`
	}{Name: "configMap0", Namespace: "klovercloud"})
	pipeline := v1.Pipeline{
		Steps: []v1.Step{v1.Step{
			Arg: variable,
		}},
	}
	testCase := TestCase{
		data:     pipeline,
		expected: map[string]string{"env1": "value1", "env2": "value2"},
	}
	service:=pipelineService{
		k8s:      &MockK8sService{},
		tekton:   nil,
		pipeline: pipeline,
	}
	service.LoadArgs(service.k8s)

	testCase.actual = service.pipeline.Steps[0].Arg.Data
	if !reflect.DeepEqual(testCase.expected, testCase.actual) {
		assert.ElementsMatch(t, testCase.expected, testCase.actual)
	}
	variable.Secrets = append(variable.Secrets, struct {
		Name      string `json:"name"`
		Namespace string `json:"namespace"`
	}{Name: "secret0", Namespace: "klovercloud"})

	pipeline.Steps[0].Arg = variable
	testCase.expected = map[string]string{"env1": "value1", "env2": "value2", "key1": "value1", "key2": "value2"}

	service.pipeline=pipeline
	service.LoadArgs(service.k8s)

	testCase.actual = service.pipeline.Steps[0].Arg.Data
	testCase.actual = pipeline.Steps[0].Arg.Data
	if !reflect.DeepEqual(testCase.expected, testCase.actual) {
		assert.ElementsMatch(t, testCase.expected, testCase.actual)
	}
}
func Test_loadEnvs(t *testing.T) {
	type TestCase struct {
		data     v1.Pipeline
		expected map[string]string
		actual   map[string]string
	}
	variable := v1.Variable{}
	variable.ConfigMaps = append(variable.ConfigMaps, struct {
		Name      string `json:"name"`
		Namespace string `json:"namespace"`
	}{Name: "configMap0", Namespace: "klovercloud"})
	pipeline := v1.Pipeline{Steps: []v1.Step{
		{
			Env: variable,
		},
	}}
	testCase := TestCase{
		data:     pipeline,
		expected: map[string]string{"env1": "value1", "env2": "value2"},
	}
	service:=pipelineService{
		k8s:      &MockK8sService{},
		tekton:   nil,
		pipeline: pipeline,
	}
	service.LoadEnvs(service.k8s)
	testCase.actual = service.pipeline.Steps[0].Env.Data
	if !reflect.DeepEqual(testCase.expected, testCase.actual) {
		assert.ElementsMatch(t, testCase.expected, testCase.actual)
	}
	variable.Secrets = append(variable.Secrets, struct {
		Name      string `json:"name"`
		Namespace string  `json:"namespace"`
	}{Name:"secret0", Namespace:"klovercloud"})
	pipeline.Steps[0].Env = variable
	testCase.expected = map[string]string{"env1": "value1", "env2": "value2", "key1": "value1", "key2": "value2"}
	service.LoadEnvs(service.k8s)
	testCase.actual = service.pipeline.Steps[0].Env.Data
	if !reflect.DeepEqual(testCase.expected, testCase.actual) {
		assert.ElementsMatch(t, testCase.expected, testCase.actual)
	}
}
