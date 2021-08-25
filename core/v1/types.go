package v1

import "github.com/klovercloud-ci/enums"

type Resource struct {
	Type enums.PIPELINE_RESOURCE_TYPE `json:"type"`
	Url string `json:"url"`
	Revision string  `json:"revision"`
}

func(resource Resource)Validate() error{

	return nil
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