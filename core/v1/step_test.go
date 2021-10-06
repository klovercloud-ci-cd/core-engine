package v1

import (
	"errors"
	"fmt"
	"github.com/klovercloud-ci/enums"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestStep_Validate(t *testing.T) {
	type TestCase struct {
		step     Step
		expected error
		actual   error
	}

	name := []string{"", "build", "Deploy", "Deploy"}
	sType := []string{"BUILD", "DEPLOY", "", "DEPLOY"}
	triger := []string{"AUTO", "", "MANUAL", "blank"}
	params1 := map[enums.PARAMS]string{"repository_type": "git", "revision": "121223234443434", "service_account": "test-sa", "images": "zeromsi2/test-dev:1.0.0,zeromsi2/test-pro:1.0.0"}
	outputs := []error{errors.New("Step name required!"), errors.New("Step triger required!"), errors.New("Step type is required!"), errors.New("Invalid triger type!")}
	params2 := map[enums.PARAMS]string{"agent": "local-agent", "name": "ubuntu", "namespace": "default", "type": "deployment", "images": "zeromsi2/test-dev:1.0.0"}

	var testCase []TestCase

	for i := 0; i < 3; i++ {
		var params = params1
		if i%2 != 0 {
			params = params2
		}
		testcase := TestCase{
			step: Step{
				Name:        name[i],
				Type:        enums.STEP_TYPE(sType[i]),
				Trigger:     enums.TRIGGER(triger[i]),
				Params:      params,
				Next:        nil,
				ArgData:     nil,
				EnvData:     nil,
				Descriptors: nil,
			},
			expected: outputs[i],
		}
		testCase = append(testCase, testcase)
	}
	for i := 0; i < 3; i++ {
		testCase[i].actual = testCase[i].step.Validate()
		if !reflect.DeepEqual(testCase[i].expected, testCase[i].actual) {
			fmt.Println(testCase[i].actual, i)
			assert.ElementsMatch(t, testCase[i].expected, testCase[i].actual)
		}
	}

	params := []map[enums.PARAMS]string{{"repository_type": "", "revision": "121223234443434", "service_account": "test-sa", "images": "zeromsi2/test-dev:1.0.0,zeromsi2/test-pro:1.0.0"}, {"repository_type": "git", "revision": "", "service_account": "test-sa", "images": "zeromsi2/test-dev:1.0.0,zeromsi2/test-pro:1.0.0"}, {"repository_type": "git", "revision": "121223234443434", "service_account": "", "images": "zeromsi2/test-dev:1.0.0,zeromsi2/test-pro:1.0.0"}, {"repository_type": "git", "revision": "121223234443434", "service_account": "test-sa", "images": ""}}

	outputs = []error{errors.New("Repository type is required!"), errors.New("Revision is required!"), errors.New("Service account is required!"), errors.New("Image is required!")}

	var testCaseForBuilld []TestCase
	for i := 0; i < 4; i++ {
		testcase := TestCase{
			step: Step{
				Name:        "build",
				Type:        "BUILD",
				Trigger:     "AUTO",
				Params:      params[i],
				Next:        nil,
				ArgData:     nil,
				EnvData:     nil,
				Descriptors: nil,
			},
			expected: outputs[i],
		}
		testCaseForBuilld = append(testCaseForBuilld, testcase)
	}
	for i := 0; i < 3; i++ {
		testCaseForBuilld[i].actual = testCaseForBuilld[i].step.Validate()
		if !reflect.DeepEqual(testCaseForBuilld[i].expected, testCaseForBuilld[i].actual) {
			assert.ElementsMatch(t, testCaseForBuilld[i].expected, testCaseForBuilld[i].actual)
		}
	}

	params = []map[enums.PARAMS]string{{"agent": "", "name": "ubuntu", "namespace": "default", "type": "deployment", "images": "zeromsi2/test-dev:1.0.0"}, {"agent": "local-agent", "name": "", "namespace": "default", "type": "deployment", "images": "zeromsi2/test-dev:1.0.0"}, {"agent": "local-agent", "name": "ubuntu", "namespace": "", "type": "deployment", "images": "zeromsi2/test-dev:1.0.0"}, {"agent": "local-agent", "name": "ubuntu", "namespace": "default", "type": "", "images": "zeromsi2/test-dev:1.0.0"}, {"agent": "local-agent", "name": "ubuntu", "namespace": "default", "type": "deployment", "images": ""}}

	outputs = []error{errors.New("Agent is required!"), errors.New("Params name is required!"), errors.New("Params namespace is required!"), errors.New("Params type is required!"), errors.New("Params image is required!")}

	var testCaseForDeploy []TestCase

	for i := 0; i < 5; i++ {
		testcase := TestCase{
			step: Step{
				Name:        "build",
				Type:        "DEPLOY",
				Trigger:     "AUTO",
				Params:      params[i],
				Next:        nil,
				ArgData:     nil,
				EnvData:     nil,
				Descriptors: nil,
			},
			expected: outputs[i],
		}
		testCaseForDeploy = append(testCaseForDeploy, testcase)
	}
	for i := 0; i < 3; i++ {
		testCaseForDeploy[i].actual = testCaseForDeploy[i].step.Validate()
		if !reflect.DeepEqual(testCaseForDeploy[i].expected, testCaseForDeploy[i].actual) {
			assert.ElementsMatch(t, testCaseForDeploy[i].expected, testCaseForDeploy[i].actual)
		}
	}
}
