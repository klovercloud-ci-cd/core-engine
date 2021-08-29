package repository

import v1 "github.com/klovercloud-ci/core/v1"

type ProcessEventRepository interface {
	Store( data v1.PipelineProcessEvent)
	GetByProcessId(processId string)map[string]interface{}
	DequeueByProcessId(processId string)map[string]interface{}
}