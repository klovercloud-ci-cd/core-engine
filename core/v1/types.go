package v1

import (
	"errors"
	"github.com/klovercloud-ci/enums"
)

type Resource struct {
	Type enums.PIPELINE_RESOURCE_TYPE `json:"type"`
	Url string `json:"url"`
	Revision string  `json:"revision"`
}

type Variable struct {
	Secrets     []struct{
		Name string `json:"name"`
		Namespace   string `json:"namespace"`
	}`json:"secrets"`
	ConfigMaps  [] struct{
		Name string `json:"name"`
		Namespace   string `json:"namespace"`
	}`json:"configMaps"`
	Data map[string]string `json:"data"`
}

type LogEventQueryOption struct {
	IndexFrom int32
	IndexTo int32
}
type PipelineApplyOption struct {
	EnablePurging bool
}
type PodListGetOption struct {
	Wait bool
	Duration int
}

func(resource Resource)Validate() error{
	if resource.Type != "" {
		if resource.Type == enums.GIT || resource.Type == enums.IMAGE {
			if resource.Url == "" {
				return errors.New("resource url is empty")
			}
			if resource.Revision == "" {
				return errors.New("resource revision is empty")
			}
			return nil
		}
		return errors.New("resource type is not match")
	}
	return errors.New("resource type is required")
}
