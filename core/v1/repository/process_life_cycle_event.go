package repository

import (
	v1 "github.com/klovercloud-ci-cd/core-engine/core/v1"
)

// ProcessLifeCycleEventRepository Process LifeCycle Event Repository operations.
type ProcessLifeCycleEventRepository interface {
	Store(data []v1.ProcessLifeCycleEvent)
	Get(count int64) []v1.ProcessLifeCycleEvent
	PullPausedAndAutoTriggerEnabledResourcesByAgentName(count int64, agent string) []v1.ProcessLifeCycleEvent
	PullNonInitializedAndAutoTriggerEnabledEventsByStepType(count int64, stepType string) []v1.ProcessLifeCycleEvent
	PullCancellingStepsByProcessStatusAndStepType(stepType string) []v1.ProcessLifeCycleEvent
}
