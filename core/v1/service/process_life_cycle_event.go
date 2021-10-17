package service

import (
	v1 "github.com/klovercloud-ci/core/v1"
)

type ProcessLifeCycleEvent interface {
	PullBuildEvents()[]v1.ProcessLifeCycleEvent
	Listen(subject v1.Subject)
}



