package v1

import (
	"errors"
	"github.com/klovercloud-ci/enums"
)

type Step struct {
	Name string `json:"name"`
	Type enums.STEP_TYPE `json:"type"`
	ServiceAccount string `json:"service_account"`
	Input Resource `json:"input"`
	Outputs []Resource `json:"outputs"`
	Arg Variable  `json:"arg"`
	Env Variable  `json:"env"`
}

func (step Step)Validate()error{
	if step.Name != ""{
		if step.Type == enums.BUILD || step.Type == enums.DEPLOY{
			if step.ServiceAccount == ""{
				return errors.New("step service account required")
			}
			err:=step.Input.Validate()
			if err!=nil{
				return err
			}

			for _,each:=range step.Outputs{
				err:=each.Validate()
				if err!=nil{
					return err
				}
			}
			return nil
		}
		if step.Type != "" {
			return errors.New("step type is required")
		}

		return errors.New("step type is not match")
	}
	return errors.New("step name required")
}