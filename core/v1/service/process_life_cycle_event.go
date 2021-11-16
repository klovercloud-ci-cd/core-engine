package service

import (
	v1 "github.com/klovercloud-ci-cd/core-engine/core/v1"
)

// ProcessLifeCycleEvent process life cycle related operations.
type ProcessLifeCycleEvent interface {
	PullBuildEvents() []v1.ProcessLifeCycleEvent
	PullBuildCancellingEvents() []v1.ProcessLifeCycleEvent
	Listen(subject v1.Subject)
}
