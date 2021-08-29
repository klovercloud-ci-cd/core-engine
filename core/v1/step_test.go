package v1

import (
	"github.com/klovercloud-ci/enums"
	"github.com/stretchr/testify/assert"
	"log"
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
	t.Run("when resource type is git", func(t *testing.T) {
		actualError := Step.Validate(step)
		assert.Equal(t, nil, actualError)
	})
	step.Input.Type = "build"
	t.Run("when resource type not match", func(t *testing.T) {
		actualError := Step.Validate(step)
		assert.Error(t, actualError)
	})
	step.Input.Type = "git"
	step.ServiceAccount = ""
	t.Run("when step service account is empty", func(t *testing.T) {
		actualError := Step.Validate(step)
		assert.Error(t, actualError)
	})
	step.ServiceAccount = "Testacc"
	step.Type = "develop"
	t.Run("when step type not match", func(t *testing.T) {
		actualError := Step.Validate(step)
		assert.Error(t, actualError)
	})
	step.Type = ""
	t.Run("when step type is empty", func(t *testing.T) {
		actualError := Step.Validate(step)
		assert.Error(t, actualError)
	})
	step.Type = enums.BUILD
	step.Input.Type = ""
	t.Run("when resource type is empty", func(t *testing.T) {
		actualError := Step.Validate(step)
		log.Print(actualError)
		assert.Error(t, actualError)
	})
}