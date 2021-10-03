package v1

import (
	"github.com/klovercloud-ci/enums"
	"time"
)

type ProcessLifeCycleEvent struct {
	ProcessId string            `bson:"process_id" json :"process_id"`
	Step string                 `bson:"step" json :"step"`
	StepType enums.STEP_TYPE    `bson:"step_type" json :"step_type"`
	Status enums.PROCESS_STATUS `bson:"status" json :"status"`
	Next []string               `bson:"next" json :"next"`
	Agent string                `bson:"agent" json :"agent"`
	Pipeline *Pipeline          `bson:"pipeline" json :"pipeline"`
	CreatedAt time.Time         `bson:"created_at" json :"created_at"`
	Trigger enums.TRIGGER       `bson:"trigger" json :"trigger"`
}