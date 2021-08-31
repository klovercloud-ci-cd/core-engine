package v1

import (
	"errors"
	"github.com/klovercloud-ci/enums"
)

type Resource struct {
	Type     enums.PIPELINE_RESOURCE_TYPE `json:"type"`
	Url      string                       `json:"url"`
	Revision string                       `json:"revision"`
}

type Variable struct {
	Secrets []struct {
		Name      string `json:"name"`
		Namespace string `json:"namespace"`
	} `json:"secrets"`
	ConfigMaps []struct {
		Name      string `json:"name"`
		Namespace string `json:"namespace"`
	} `json:"configMaps"`
	Data map[string]string `json:"data"`
}

type LogEventQueryOption struct {
	Pagination struct {
		Page  int64
		Limit int64
	}
	Step string
}
type PipelineApplyOption struct {
	EnablePurging bool
}
type PodListGetOption struct {
	Wait     bool
	Duration int
}
type Listener struct {
	Name,Namespace,ContainerName,Step,ProcessId,Log string
	StepType enums.STEP_TYPE
	EventData map[string]interface{}
	ProcessLabel map[string]string
}
func (resource Resource) Validate() error {
	if resource.Type!=enums.IMAGE &&  resource.Type!=enums.GIT{
		return errors.New("Invalid resource type!")
	}
	if resource.Url == "" {
		return errors.New("Resource url is required!")
	}
	if resource.Revision == "" {
		return errors.New("Resource revision is required!")
	}
	return nil
}
