package v1

import (
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	"reflect"
	"testing"
)

func Test_setEnvs(t *testing.T) {
	type TestCase struct {
		configMaps [] v1.ConfigMap
		secrets [] v1.Secret
		expected map[string]string
		actual map[string]string
	}

	testCase:= TestCase{
		configMaps:     InitConfigMaps(),
		secrets: InitSecrets(),
		expected: map[string]string{"key1":"value1","key2":"value2","env1":"value1","env2":"value2"},
	}

	step := &Step{
		Env: Variable{
			Secrets: []struct {
				Name      string `json:"name"`
				Namespace string `json:"namespace"`
			}{{Name: "secret0", Namespace:"klovercloud" }},
			ConfigMaps: []struct {
				Name      string `json:"name"`
				Namespace string `json:"namespace"`
			}{{Name: "configMap0", Namespace: "klovercloud"}},
		},
	}
	step.setEnvs(&MockK8sResource{})
	testCase.actual=step.Env.Data

	if !reflect.DeepEqual(testCase.expected, testCase.actual){
		assert.ElementsMatch(t, testCase.expected, testCase.actual)
	}
}

func Test_setArgs(t *testing.T) {
	type TestCase struct {
		configMaps [] v1.ConfigMap
		secrets [] v1.Secret
		expected map[string]string
		actual map[string]string
	}

	testCase:= TestCase{
		configMaps:     InitConfigMaps(),
		secrets: InitSecrets(),
		expected: map[string]string{"key1":"value1","key2":"value2","env1":"value1","env2":"value2"},
		actual:   nil,
	}

	step := &Step{
		Arg: Variable{
			Secrets: []struct {
				Name      string `json:"name"`
				Namespace string `json:"namespace"`
			}{{Name: "secret0", Namespace:"klovercloud" }},
			ConfigMaps: []struct {
				Name      string `json:"name"`
				Namespace string `json:"namespace"`
			}{{Name: "configMap0", Namespace: "klovercloud"}},
		},
	}
	step.setArgs(&MockK8sResource{})
	testCase.actual=step.Arg.Data

	if !reflect.DeepEqual(testCase.expected, testCase.actual){
		assert.ElementsMatch(t, testCase.expected, testCase.actual)
	}
}

func Test_setInput(t *testing.T) {
	type TestCase struct {
		url,revision string
		step_type STEP_TYPE
		expected Resource
		actual Resource
	}
	step := Step{
		Type: BUILD,
	}
	testCase:=TestCase{
		url:       "github.com/abc",
		revision:  "1222",
		step_type: BUILD,
		expected:  Resource{
			Url:      "github.com/abc",
			Revision: "1222",
		},
	}

	step.setInput("github.com/abc","1222")
	testCase.actual=step.Input
	if !reflect.DeepEqual(testCase.expected, testCase.actual){
		assert.ElementsMatch(t, testCase.expected, testCase.actual)
	}
	testCase=TestCase{
		url:       "github.com/abc",
		revision:  "1222",
		step_type: DEPLOY,
		expected:  Resource{},
	}
	step.setInput("github.com/abc","1222")
	if !reflect.DeepEqual(testCase.expected, testCase.actual){
		assert.ElementsMatch(t, testCase.expected, testCase.actual)
	}
}