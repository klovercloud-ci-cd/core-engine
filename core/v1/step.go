package v1

import (
	"errors"
	"github.com/klovercloud-ci/enums"
)

type Step struct {
	Name string `json:"name" yaml:"name"`
	Type enums.STEP_TYPE `json:"type" yaml:"type"`
	ServiceAccount string `json:"service_account" yaml:"service_account"`
	Input Resource `json:"input"  yaml:"input"`
	Outputs []Resource `json:"outputs"  yaml:"outputs"`
	Arg Variable  `json:"arg"  yaml:"arg"`
	Env Variable  `json:"env"  yaml:"env"`
}

func (step Step)Validate()error{
	if step.Name == ""{
		return errors.New("Step name required!")
	}
	if step.Type!=enums.BUILD && step.Type!=enums.DEPLOY{
		return errors.New("Invalid step type!")
	}
	if step.Type==enums.BUILD{
		err:=step.Input.Validate()
		if err!=nil{
			return err
		}
	}

	for _,each:=range step.Outputs{
		err:=each.Validate()
		if err!=nil{
			return err
		}
	}
	if step.Type==enums.BUILD {
		if step.ServiceAccount == "" {
			return errors.New("Service account required!")
		}
	}
return nil
}