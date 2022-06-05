package service

import (
	v1 "github.com/klovercloud-ci-cd/core-engine/core/v1"
)

// ProcessEvent Process event operations.
type ProcessEvent interface {
	Store(data v1.ProcessEvent)
	GetByCompanyId(companyId string) map[string]interface{}
	DequeueByCompanyId(companyId string) map[string]interface{}
	Listen(subject v1.Subject)
}
