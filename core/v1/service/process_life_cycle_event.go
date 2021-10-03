package service

import (
	v1 "github.com/klovercloud-ci/core/v1"
)

type ProcessLifeCycleEvent interface {
	Store( events []v1.ProcessLifeCycleEvent)
	GetByProcessId(processId string)[]v1.ProcessLifeCycleEvent
	Listen(subject v1.Subject)
}



