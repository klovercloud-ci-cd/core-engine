package repository

import v1 "github.com/klovercloud-ci-cd/core-engine/core/v1"

// ProcessEventRepository Process Event Repository operations.
type ProcessEventRepository interface {
	Store(data v1.ProcessEvent)
	GetByCompanyId(companyId string) map[string]interface{}
	DequeueByCompanyId(companyId string) map[string]interface{}
}
