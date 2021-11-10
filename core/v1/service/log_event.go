package service

import (
	v1 "github.com/klovercloud-ci-cd/core-engine/core/v1"
)

// LogEvent LogEvent operations.
type LogEvent interface {
	Store(log v1.LogEvent)
	GetByProcessId(processId string, option v1.LogEventQueryOption) ([]string, int64)
	Listen(subject v1.Subject)
}
