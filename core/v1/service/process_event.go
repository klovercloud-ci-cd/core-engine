package service

import (
	v1 "github.com/klovercloud-ci/core/v1"
)

// ProcessEvent Process event operations.
type ProcessEvent interface {
	Store(data v1.ProcessEvent)
	GetByProcessId(processId string) map[string]interface{}
	DequeueByProcessId(processId string) map[string]interface{}
	Listen(subject v1.Subject)
}
