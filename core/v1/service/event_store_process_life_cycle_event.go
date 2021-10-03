package service

import (
	v1 "github.com/klovercloud-ci/core/v1"
)
type EventStoreProcessLifeCycleEvent interface {
	Listen(subject v1.Subject)
}
