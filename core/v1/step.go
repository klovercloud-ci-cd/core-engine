package v1

import (
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

