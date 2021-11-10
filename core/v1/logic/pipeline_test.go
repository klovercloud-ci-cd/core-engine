package logic

import (
	"fmt"
	v1 "github.com/klovercloud-ci-cd/core-engine/core/v1"
	"github.com/klovercloud-ci-cd/core-engine/enums"
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
	pipeline := v1.Pipeline{
		Steps: []v1.Step{{
			Params: map[enums.PARAMS]string{"args_from_configmaps": "klovercloud/configMap1"},
		}},
	}
	testCase := TestCase{
		data:     pipeline,
		expected: map[string]string{"env1": "value1", "env2": "value2"},
	}
	service := pipelineService{
		k8s:    &mockK8sService{},
		tekton: nil,
	}
	service.LoadArgs(pipeline)
	testCase.actual = service.pipeline.Steps[0].ArgData
	if !reflect.DeepEqual(testCase.expected, testCase.actual) {
		fmt.Println(testCase.actual)
		assert.ElementsMatch(t, testCase.expected, testCase.actual)
	}
}

func Test_loadEnvs(t *testing.T) {
	type TestCase struct {
		data     v1.Pipeline
		expected map[string]string
		actual   map[string]string
	}
	pipeline := v1.Pipeline{Steps: []v1.Step{
		{
			Params: map[enums.PARAMS]string{"envs_from_configmaps": "klovercloud/configMap1"},
		},
	}}
	testCase := TestCase{
		data:     pipeline,
		expected: map[string]string{"env1": "value1", "env2": "value2"},
	}
	service := pipelineService{
		k8s:    &mockK8sService{},
		tekton: nil,
	}
	service.LoadEnvs(pipeline)
	testCase.actual = service.pipeline.Steps[0].EnvData
	if !reflect.DeepEqual(testCase.expected, testCase.actual) {
		fmt.Println(testCase.actual)
		assert.ElementsMatch(t, testCase.expected, testCase.actual)
	}
}

func TestPipelineService_SetInputResource(t *testing.T) {
	type TestCase struct {
		data     v1.Pipeline
		expected map[enums.PARAMS]string
		actual   map[enums.PARAMS]string
	}
	pipeline := v1.Pipeline{Steps: []v1.Step{
		{
			Name:    "test",
			Type:    "BUILD",
			Trigger: "AUTO",
			Params:  map[enums.PARAMS]string{"envs_from_configmaps": "klovercloud/configMap1"},
		},
	},
	}
	testCase := TestCase{
		data:     pipeline,
		expected: map[enums.PARAMS]string{"envs_from_configmaps": "klovercloud/configMap1", "revision": "123456", "url": "www.example.com"},
	}
	service := pipelineService{
		k8s:    &mockK8sService{},
		tekton: nil,
	}
	service.SetInputResource("www.example.com", "123456", testCase.data)
	testCase.actual = service.pipeline.Steps[0].Params
	if !reflect.DeepEqual(testCase.expected, testCase.actual) {
		fmt.Println(testCase.actual)
		assert.ElementsMatch(t, testCase.expected, testCase.actual)
	}
}
func TestPipelineService_Build(t *testing.T) {
	type TestCase struct {
		data     v1.Pipeline
		expected v1.Step
		actual   v1.Step
	}
	pipeline := v1.Pipeline{Steps: []v1.Step{
		{
			Name:    "test",
			Type:    "BUILD",
			Trigger: "AUTO",
			Params:  map[enums.PARAMS]string{"envs_from_configmaps": "klovercloud/configMap1", "args_from_configmaps": "klovercloud/configMap1"},
		},
	},
	}
	testCase := TestCase{
		data: pipeline,
		expected: v1.Step{
			Name:    "test",
			Type:    "BUILD",
			Trigger: "AUTO",
			Params:  map[enums.PARAMS]string{"envs_from_configmaps": "klovercloud/configMap1", "args_from_configmaps": "klovercloud/configMap1", "revision": "123456", "url": "www.example.com"},
			Next:    []string{},
			ArgData: map[string]string{"env1": "value1", "env2": "value2"},
			EnvData: map[string]string{"env1": "value1", "env2": "value2"},
		},
	}
	service := pipelineService{
		k8s:    &mockK8sService{},
		tekton: nil,
	}
	service.Build("www.example.com", "123456", testCase.data)
	testCase.actual = service.pipeline.Steps[0]
	//if !reflect.DeepEqual(testCase.expected, testCase.actual) {
	//	fmt.Println(testCase.actual, testCase.expected)
	//	assert.ElementsMatch(t, testCase.expected, testCase.actual)
	//}
}
