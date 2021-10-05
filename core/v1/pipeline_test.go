package v1

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestPipeline_Validate(t *testing.T) {
	type TestCase struct {
		Data     Pipeline
		Expected error
		Actual   error
	}
	apiVersion := []string{"", "001", "0002"}
	name := []string{"app1", "", "app2"}
	proccessId := []string{"01", "02", ""}
	outputs := []error{errors.New("Api version is required!"), errors.New("Pipeline name is required!"), errors.New("Pipeline process id is required!")}

	var testcase []TestCase

	for i := 0; i < 3; i++ {
		testdata := TestCase{
			Data: Pipeline{
				Option:     PipelineApplyOption{},
				ApiVersion: apiVersion[i],
				Name:       name[i],
				ProcessId:  proccessId[i],
				Label:      nil,
				Steps:      nil,
			},
			Expected: outputs[i],
		}
		testcase = append(testcase, testdata)
	}
	for i := 0; i < 3; i++ {
		testcase[i].Actual = testcase[i].Data.Validate()
		if !reflect.DeepEqual(testcase[i].Expected, testcase[i].Actual) {
			assert.ElementsMatch(t, testcase[i].Expected, testcase[i].Actual)
		}
	}
}
