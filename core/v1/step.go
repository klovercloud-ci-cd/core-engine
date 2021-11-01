package v1

import (
	"errors"
	"github.com/klovercloud-ci/enums"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

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

func (step Step) Validate() error {
	if step.Name == "" {
		return errors.New("Step name required!")
	}
	if step.Type == enums.BUILD {
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
	} else if step.Type == enums.DEPLOY {
		if step.Params[enums.AGENT] == "" {
			return errors.New("Agent is required!")
		}
		if step.Params[enums.NAME] == "" {
			return errors.New("Params name is required!")
		}
		if step.Params[enums.NAMESPACE] == "" {

			return errors.New("Params namespace is required!")
		}
		if step.Params[enums.TYPE] == "" {
			return errors.New("Params type is required!")
		}
		if step.Params[enums.IMAGES] == "" {
			return errors.New("Params image is required!")
		}
	} else if step.Type == "" {
		return errors.New("Step type is required!")
	} else {
		return errors.New("Invalid step type!")
	}
	if step.Trigger != enums.AUTO && step.Trigger != enums.MANUAL {
		return errors.New("Invalid trigger type!")
	}
	return nil
}
