package logic

import (
	v1 "github.com/klovercloud-ci/core/v1"
	"github.com/klovercloud-ci/enums"
	"github.com/stretchr/testify/assert"
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

//
//import (
//	v1 "github.com/klovercloud-ci/core/v1"
//	"github.com/klovercloud-ci/enums"
//	"github.com/stretchr/testify/assert"
//	corev1 "k8s.io/api/core/v1"
//	"reflect"
//	"testing"
//)
//
//
//func Test_SetEnvs(t *testing.T) {
//	type TestCase struct {
//		configMaps [] corev1.ConfigMap
//		secrets [] corev1.Secret
//		expected map[string]string
//		actual map[string]string
//	}
//
//	testCase:= TestCase{
//		configMaps:     InitConfigMaps(),
//		secrets: InitSecrets(),
//		expected: map[string]string{"key1":"value1","key2":"value2","env1":"value1","env2":"value2"},
//	}
//
//	step := &v1.Step{
//		Env: v1.Variable{
//			Secrets: []struct {
//				Name      string `json:"name" yaml:"name"`
//				Namespace string `json:"namespace" yaml:"namespace"`
//			}([]struct {
//				Name      string `json:"name"`
//				Namespace string `json:"namespace"`
//			}{{Name: "secret0", Namespace: "klovercloud"}}),
//			ConfigMaps: []struct {
//				Name      string `json:"name" yaml:"name"`
//				Namespace string `json:"namespace" yaml:"namespace"`
//			}([]struct {
//				Name      string `json:"name"`
//				Namespace string `json:"namespace"`
//			}{{Name: "configMap0", Namespace: "klovercloud"}}),
//		},
//	}
//	service:=stepService{
//		step: *step,
//	}
//	service.SetEnvs(&mockK8sService{})
//	testCase.actual=service.step.Env.Data
//	if !reflect.DeepEqual(testCase.expected, testCase.actual){
//		assert.ElementsMatch(t, testCase.expected, testCase.actual)
//	}
//}
//
//func Test_SetArgs(t *testing.T) {
//	type TestCase struct {
//		configMaps [] corev1.ConfigMap
//		secrets [] corev1.Secret
//		expected map[string]string
//		actual map[string]string
//	}
//
//	testCase:= TestCase{
//		configMaps:     InitConfigMaps(),
//		secrets: InitSecrets(),
//		expected: map[string]string{"key1":"value1","key2":"value2","env1":"value1","env2":"value2"},
//		actual:   nil,
//	}
//
//	step := &v1.Step{
//		Arg: v1.Variable{
//			Secrets: []struct {
//				Name      string `json:"name" yaml:"name"`
//				Namespace string `json:"namespace" yaml:"namespace"`
//			}([]struct {
//				Name      string `json:"name"`
//				Namespace string `json:"namespace"`
//			}{{Name: "secret0", Namespace: "klovercloud"}}),
//			ConfigMaps: []struct {
//				Name      string `json:"name" yaml:"name"`
//				Namespace string `json:"namespace" yaml:"namespace"`
//			}([]struct {
//				Name      string `json:"name"`
//				Namespace string `json:"namespace"`
//			}{{Name: "configMap0", Namespace: "klovercloud"}}),
//		},
//	}
//	service:=stepService{
//		step: *step,
//	}
//	service.SetArgs(&mockK8sService{})
//	testCase.actual=service.step.Arg.Data
//	if !reflect.DeepEqual(testCase.expected, testCase.actual){
//		assert.ElementsMatch(t, testCase.expected, testCase.actual)
//	}
//}
//
//func Test_SetInput(t *testing.T) {
//	type TestCase struct {
//		url,revision string
//		step_type enums.STEP_TYPE
//		expected v1.Resource
//		actual v1.Resource
//	}
//	step := v1.Step{
//		Type: enums.BUILD,
//	}
//	testCase:=TestCase{
//		url:       "github.com/abc",
//		revision:  "1222",
//		step_type: enums.BUILD,
//		expected:  v1.Resource{
//			Url:      "github.com/abc",
//			Revision: "1222",
//		},
//	}
//
//	service:=stepService{
//		step: step,
//	}
//	service.SetInput("github.com/abc","1222")
//	testCase.actual=service.step.Input
//
//	if !reflect.DeepEqual(testCase.expected, testCase.actual){
//		assert.ElementsMatch(t, testCase.expected, testCase.actual)
//	}
//	testCase=TestCase{
//		url:       "github.com/abc",
//		revision:  "1222",
//		step_type: enums.DEPLOY,
//		expected:  v1.Resource{},
//	}
//	service.SetInput("github.com/abc","1222")
//	if !reflect.DeepEqual(testCase.expected, testCase.actual){
//		assert.ElementsMatch(t, testCase.expected, testCase.actual)
//	}
//}
