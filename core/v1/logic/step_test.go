package logic

import (
	"fmt"
	v1 "github.com/klovercloud-ci-cd/klovercloud-ci-core/core/v1"
	"github.com/klovercloud-ci-cd/klovercloud-ci-core/enums"
	"github.com/stretchr/testify/assert"
	"log"
	"reflect"
	"testing"
)

func TestStepService_SetInput(t *testing.T) {
	type TestCase struct {
		data     v1.Step
		expected map[enums.PARAMS]string
		actual   map[enums.PARAMS]string
	}

	testCase := TestCase{
		expected: map[enums.PARAMS]string{"repository_type": "git", "revision": "123456", "service_account": "test-sa", "images": "zeromsi2/test-dev:1.0.0,zeromsi2/test-pro:1.0.0", "url": "www.example.com"},
	}
	step := &v1.Step{
		Name:        "build",
		Type:        "BUILD",
		Trigger:     "AUTO",
		Params:      map[enums.PARAMS]string{"repository_type": "git", "revision": "121223234443434", "service_account": "test-sa", "images": "zeromsi2/test-dev:1.0.0,zeromsi2/test-pro:1.0.0"},
		Next:        nil,
		ArgData:     nil,
		EnvData:     nil,
		Descriptors: nil,
	}
	service := stepService{
		step: *step,
	}
	service.SetInput("www.example.com", "123456")
	testCase.actual = service.step.Params
	if !reflect.DeepEqual(testCase.expected, testCase.actual) {
		assert.ElementsMatch(t, testCase.expected, testCase.actual)
	}
}

func TestStepService_SetArgs(t *testing.T) {
	type TestCase struct {
		data     v1.Step
		expected map[string]string
		actual   map[string]string
	}
	testCase := TestCase{
		expected: map[string]string{"env1": "value1", "env2": "value2"},
	}
	step := &v1.Step{
		Name:        "build",
		Type:        "BUILD",
		Trigger:     enums.AUTO,
		Params:      map[enums.PARAMS]string{"args_from_configmaps": "klovercloud/configMap1"},
		Next:        nil,
		ArgData:     nil,
		EnvData:     nil,
		Descriptors: nil,
	}
	service := stepService{
		step: *step,
	}
	service.SetArgs(NewMockK8sService(nil))
	testCase.actual = service.step.ArgData
	if !reflect.DeepEqual(testCase.expected, testCase.actual) {
		log.Println(testCase.expected, testCase.actual)
		assert.ElementsMatch(t, testCase.expected, testCase.actual)
	}
}

func TestStepService_SetEnvs(t *testing.T) {
	type TestCase struct {
		data     v1.Step
		expected map[string]string
		actual   map[string]string
	}
	testCase := TestCase{
		expected: map[string]string{"env1": "value1", "env2": "value2"},
	}
	step := &v1.Step{
		Name:        "build",
		Type:        "BUILD",
		Trigger:     "AUTO",
		Params:      map[enums.PARAMS]string{"envs_from_configmaps": "klovercloud/configMap1"},
		Next:        nil,
		ArgData:     nil,
		EnvData:     nil,
		Descriptors: nil,
	}
	service := stepService{
		step: *step,
	}
	service.SetEnvs(NewMockK8sService(nil))
	testCase.actual = service.step.EnvData
	if !reflect.DeepEqual(testCase.expected, testCase.actual) {
		fmt.Println(testCase.expected, testCase.actual)
		assert.ElementsMatch(t, testCase.expected, testCase.actual)
	}
}
