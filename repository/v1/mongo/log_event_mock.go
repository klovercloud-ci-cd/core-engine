package mongo

import (
	"context"
	v1 "github.com/klovercloud-ci-cd/core-engine/core/v1"
	"github.com/klovercloud-ci-cd/core-engine/core/v1/repository"
	"time"
)

var data []v1.LogEvent

// InitLogEventData Return mock Log event list
func InitLogEventData() []v1.LogEvent {
	data = []v1.LogEvent{
		{
			ProcessId: "1",
			Log:       "Initializing pod",
			Step:      "buildImage",
			CreatedAt: time.Time{},
		},
		{
			ProcessId: "1",
			Log:       "Pulling Image",
			Step:      "buildImage",
			CreatedAt: time.Time{},
		},
		{
			ProcessId: "2",
			Log:       "Failed to initialize pod",
			Step:      "buildImage",
			CreatedAt: time.Time{},
		},
		{
			ProcessId: "2",
			Log:       "Initializing pod",
			Step:      "deployImage",
			CreatedAt: time.Time{},
		},
		{
			ProcessId: "2",
			Log:       "Pulling Image",
			Step:      "deployImage",
			CreatedAt: time.Time{},
		},
	}

	return data
}

// NewMockLogEventRepository returns LogEventRepository type object
func NewMockLogEventRepository() repository.LogEventRepository {
	manager := GetMockDmManager()
	manager.Db.Drop(context.Background())
	return &logEventRepository{
		manager: GetMockDmManager(),
		timeout: 3000,
	}

}
