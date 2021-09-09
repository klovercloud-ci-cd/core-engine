package v1

import (
	"github.com/klovercloud-ci/enums"
	"github.com/stretchr/testify/assert"
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
			testcases[i].actual =err.Error()
		}else{
			testcases[i].actual=""
		}
		if !reflect.DeepEqual(testcases[i].expected, testcases[i].actual) {
			assert.ElementsMatch(t, testcases[i].expected, testcases[i].actual)
		}
	}

}