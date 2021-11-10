package inmemory

import (
	v1 "github.com/klovercloud-ci-cd/klovercloud-ci-core/core/v1"
	"github.com/klovercloud-ci-cd/klovercloud-ci-core/core/v1/repository"
)

type processLifeCycleRepository struct {
}

func (p processLifeCycleRepository) PullNonInitializedAndAutoTriggerEnabledEventsByStepType(count int64, stepType string) []v1.ProcessLifeCycleEvent {
	return nil
}

func (p processLifeCycleRepository) PullPausedAndAutoTriggerEnabledResourcesByAgentName(count int64, agent string) []v1.ProcessLifeCycleEvent {
	return nil
}

func (p processLifeCycleRepository) Get(count int64) []v1.ProcessLifeCycleEvent {
	return nil
}

func (p processLifeCycleRepository) Store(events []v1.ProcessLifeCycleEvent) {

}

func (p processLifeCycleRepository) updateStatus(data v1.ProcessLifeCycleEvent, status string) error {
	return nil
}
func (p processLifeCycleRepository) update(data v1.ProcessLifeCycleEvent) error {
	return nil
}
func (p processLifeCycleRepository) GetByProcessIdAndStep(processId, step string) *v1.ProcessLifeCycleEvent {
	return nil
}

func (p processLifeCycleRepository) GetByProcessId(processId string) []v1.ProcessLifeCycleEvent {
	return nil
}

// NewProcessLifeCycleRepository returns ProcessLifeCycleEventRepository type object
func NewProcessLifeCycleRepository() repository.ProcessLifeCycleEventRepository {
	return &processLifeCycleRepository{}

}
