package v1

import (
	"github.com/klovercloud-ci/enums"
	"github.com/stretchr/testify/assert"
	"log"
	"reflect"
	"testing"
)

func Test_Validate(t *testing.T) {
	type TestCase struct {
		data Step
		expected string
		actual string
	}
	step := Step{
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
	}
	testcases:=[]TestCase{}
	testcases= append(testcases, TestCase{
		data:     step,
		expected: "",
	})

	step.Input.Type = ""
	testcases= append(testcases, TestCase{
		data:     step,
		expected: "Invalid resource type!",
	})

	step.Input.Type = "git"
	step.ServiceAccount=""
	testcases= append(testcases, TestCase{
		data:     step,
		expected: "Service account required!",
	})
	for i,_:=range testcases{
		err:=Step.Validate(testcases[i].data)
		if err!=nil{
			log.Print(err.Error())
			testcases[i].actual =err.Error()
		}else{
			testcases[i].actual=""
		}
		if !reflect.DeepEqual(testcases[i].expected, testcases[i].actual) {
			assert.ElementsMatch(t, testcases[i].expected, testcases[i].actual)
		}
	}

	//t.Run("when resource type is git", func(t *testing.T) {
	//	actualError := Step.Validate(step)
	//	assert.Equal(t, nil, actualError)
	//})
	//step.Input.Type = "build"
	//t.Run("when resource type not match", func(t *testing.T) {
	//	actualError := Step.Validate(step)
	//	assert.Error(t, actualError)
	//})
	//step.Input.Type = "git"
	//step.ServiceAccount = ""
	//t.Run("when step service account is empty", func(t *testing.T) {
	//	actualError := Step.Validate(step)
	//	assert.Error(t, actualError)
	//})
	//step.ServiceAccount = "Testacc"
	//step.Type = "develop"
	//t.Run("when step type not match", func(t *testing.T) {
	//	actualError := Step.Validate(step)
	//	assert.Error(t, actualError)
	//})
	//step.Type = ""
	//t.Run("when step type is empty", func(t *testing.T) {
	//	actualError := Step.Validate(step)
	//	assert.Error(t, actualError)
	//})
	//step.Type = enums.BUILD
	//step.Input.Type = ""
	//t.Run("when resource type is empty", func(t *testing.T) {
	//	actualError := Step.Validate(step)
	//	log.Print(actualError)
	//	assert.Error(t, actualError)
	//})
}