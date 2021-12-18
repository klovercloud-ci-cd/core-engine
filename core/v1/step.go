package v1

import (
	"errors"
	"github.com/klovercloud-ci-cd/core-engine/enums"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// Step pipeline step.
type Step struct {
	Name        string                       `json:"name" yaml:"name"`
	Type        enums.STEP_TYPE              `json:"type" yaml:"type"`
	Trigger     enums.TRIGGER                `json:"trigger" yaml:"trigger"`
	Params      map[enums.PARAMS]string      `json:"params" yaml:"params"`
	Next        []string                     `json:"next" yaml:"next"`
	ArgData     map[string]string            `json:"arg_data"  yaml:"arg_data"`
	EnvData     map[string]string            `json:"env_data"  yaml:"env_data"`
	Descriptors *[]unstructured.Unstructured `json:"descriptors" yaml:"descriptors"`
}

// Validate validates pipeline steps.
func (step Step) Validate() error {
	if step.Name == "" {
		return errors.New("Step name required!")
	}
	if step.Type == enums.BUILD {
		err := step.validateBuildStep()
		if err != nil {
			return err
		}
	} else if step.Type == enums.DEPLOY {
		err := step.validateDeployStep()
		if err != nil {
			return err
		}
	} else if step.Type == enums.INTERMEDIARY {
		err := step.validateIntermediaryStep()
		if err != nil {
			return err
		}
	} else if step.Type == "" {
		return errors.New("Step type is required!")
	} else {
		return errors.New("Invalid step type!"+string(step.Type))
	}
	if step.Trigger == "" {
		return errors.New("Step triger required!")
	} else if step.Trigger != enums.AUTO && step.Trigger != enums.MANUAL {
		return errors.New("Invalid trigger type!")
	}
	return nil
}

func (step Step) validateDeployStep() error {
	if step.Params[enums.AGENT] == "" {
		return errors.New("Agent is required!")
	}
	if step.Params[enums.RESOURCE_NAME] == "" {
		return errors.New("Params name is required!")
	}
	if step.Params[enums.RESOURCE_NAMESPACE] == "" {

		return errors.New("Params namespace is required!")
	}
	if step.Params[enums.TYPE] == "" {
		return errors.New("Params type is required!")
	}
	if step.Params[enums.IMAGES] == "" {
		return errors.New("Params image is required!")
	}
	return nil
}

func (step Step) validateBuildStep() error {
	if step.Params[enums.REPOSITORY_TYPE] == "" {
		return errors.New("Repository type is required!")
	}
	if step.Params[enums.REVISION] == "" {
		return errors.New("Revision is required!")
	}
	if step.Params[enums.SERVICE_ACCOUNT] == "" {
		return errors.New("Service account is required!")
	}
	if step.Params[enums.IMAGES] == "" {
		return errors.New("Image is required!")
	}
	return nil
}

func (step Step) validateIntermediaryStep() error {
	if step.Params[enums.SERVICE_ACCOUNT] == "" {
		return errors.New("Service account is required!")
	}
	if step.Params[enums.IMAGES] == "" {
		return errors.New("Image is required!")
	}
	return nil
}
