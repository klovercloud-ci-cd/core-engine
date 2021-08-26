package v1

import "errors"

type Pipeline struct {
	Option     PipelineApplyOption
	ApiVersion string            `json:"api_version"`
	Name       string            `json:"name"`
	ProcessId  string            `json:"process_id"`
	Label      map[string]string `json:"label"`
	Steps      [] Step           `json:"steps"`
}
func(pipeline Pipeline)Validate()error{
	if pipeline.ApiVersion == "" {
		return errors.New("pipeline api version is required")
	}
	if pipeline.Name == "" {
		return errors.New("pipeline name is required")
	}
	if pipeline.ProcessId == "" {
		return errors.New("pipeline process id is required")
	}
	if pipeline.Label == nil {
		return errors.New("pipeline label is required")
	}
	for _,each:=range pipeline.Steps{
		err:=each.Validate()
		if err!=nil{
			return err
		}
	}
	return nil
}