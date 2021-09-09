package v1

import (
	"errors"
	"github.com/klovercloud-ci/enums"
)

type Resource struct {
	Type     enums.PIPELINE_RESOURCE_TYPE `json:"type" yaml:"type"`
	Url      string                       `json:"url"  yaml:"url"`
	Revision string                       `json:"revision"  yaml:"revision"`
	DeploymentResource *DeploymentResource  `json:"deployment_resource"  yaml:"deployment_resource"`
}

type DeploymentResource struct {
	ProcessId string `json:"processId" yaml:"processId"`
	Agent string `json:"agent" yaml:"agent"`
	Name string                  `json:"name" yaml:"name"`
	Namespace string             `json:"namespace" yaml:"namespace"`
	Replica int32                `json:"replica" yaml:"replica"`
	Images [] struct {
		ImageIndex int `json:"image_index" yaml:"image_index"`
		Image      string `json:"image" yaml:"image"`
	}`json:"images" yaml:"images"`
}

type Variable struct {
	Secrets []struct {
		Name      string `json:"name" yaml:"name"`
		Namespace string `json:"namespace" yaml:"namespace"`
	} `json:"secrets" yaml:"secrets"`
	ConfigMaps []struct {
		Name      string `json:"name" yaml:"name"`
		Namespace string `json:"namespace" yaml:"namespace"`
	} `json:"configMaps"  yaml:"configMaps"`
	Data map[string]string `json:"data"  yaml:"data"`
}

type Agent struct {
	Url, Token string
}
type LogEventQueryOption struct {
	Pagination struct {
		Page  int64
		Limit int64
	}
	Step string
}
type PipelineApplyOption struct {
	Purging enums.PIPELINE_PURGING
}
type PodListGetOption struct {
	Wait     bool
	Duration int
}
type Subject struct {
	Step,Log string
	StepType enums.STEP_TYPE
	EventData map[string]interface{}
	ProcessLabel map[string]string
	Pipeline Pipeline
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
