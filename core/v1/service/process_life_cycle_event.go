package service

import (
	v1 "github.com/klovercloud-ci-cd/klovercloud-ci-core/core/v1"
)

// ProcessLifeCycleEvent process life cycle related operations.
type ProcessLifeCycleEvent interface {
	PullBuildEvents() []v1.ProcessLifeCycleEvent
	Listen(subject v1.Subject)
}
