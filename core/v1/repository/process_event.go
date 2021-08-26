package repository

import v1 "github.com/klovercloud-ci/core/v1"

type ProcessEventRepository interface {
	Store( data v1.PipelineProcessStatus)
	GetByProcessId(processId string)map[string]interface{}
}