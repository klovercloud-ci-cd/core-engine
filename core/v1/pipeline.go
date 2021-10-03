package v1

import "errors"

type Pipeline struct {
	Option     PipelineApplyOption `json:"option" yaml:"option"`
	ApiVersion string              `json:"api_version" yaml:"api_version"`
	Name       string              `json:"name"  yaml:"name"`
	ProcessId  string              `json:"process_id" yaml:"process_id"`
	Label      map[string]string   `json:"label" yaml:"label"`
	Steps []Step                   `json:"steps" yaml:"steps"`
}

func(pipeline Pipeline)Validate()error{
	if pipeline.ApiVersion == "" {
		return errors.New("Api version is required!")
	}
	if pipeline.Name == "" {
		return errors.New("Pipeline name is required!")
	}
	if pipeline.ProcessId == "" {
		return errors.New("Pipeline process id is required!")
	}

	for _,each:=range pipeline.Steps{
		err:=each.Validate()
		if err!=nil{
			return err
		}
	}
	return nil
}