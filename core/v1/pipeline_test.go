package v1

import (
	"github.com/klovercloud-ci/enums"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPipeline_Validate(t *testing.T) {
	type TestCase struct {
		data Pipeline
		expected string
		actual string
	}
	pipeline := Pipeline{
		ApiVersion: "default",
		Name:       "test",
		BuildId:    "0125632",
		Label:       map[string]string{"env1": "value1", "env2": "value2"},
		Steps:      []Step{{
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
		}},
	}
	t.Run("when pipeline is validated", func(t *testing.T) {
		actualError := Pipeline.Validate(pipeline)
		assert.Equal(t, nil, actualError)
	})
	pipeline.ApiVersion = ""
	t.Run("when pipeline api version is empty", func(t *testing.T) {
		actualError := Pipeline.Validate(pipeline)
		assert.Error(t, actualError)
	})
	pipeline.ApiVersion = "default"
	pipeline.Name = ""
	t.Run("when pipeline name is empty", func(t *testing.T) {
		actualError := Pipeline.Validate(pipeline)
		assert.Error(t, actualError)
	})
	pipeline.Name = "test"
	pipeline.BuildId = ""
	t.Run("when pipeline build id is empty", func(t *testing.T) {
		actualError := Pipeline.Validate(pipeline)
		assert.Error(t, actualError)
	})
	pipeline.BuildId = "0125632"
	pipeline.Label = nil
	t.Run("when pipeline label is nil", func(t *testing.T) {
		actualError := Pipeline.Validate(pipeline)
		assert.Error(t, actualError)
	})
}