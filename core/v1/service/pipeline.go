package service

import (
	v1 "github.com/klovercloud-ci-cd/core-engine/core/v1"
	"github.com/klovercloud-ci-cd/core-engine/enums"
)

// Pipeline Pipeline operations.
type Pipeline interface {
	BuildProcessLifeCycleEvents(pipeline v1.Pipeline) error
	LoadArgs(pipeline v1.Pipeline)
	LoadEnvs(pipeline v1.Pipeline)
	SetInputResource(pipeline v1.Pipeline)
	Build(pipeline v1.Pipeline)
	PostOperations(step string, stepType enums.STEP_TYPE, pipeline v1.Pipeline, claim int)
	GetLogsByProcessId(processId string, option v1.LogEventQueryOption) ([]string, int64)
	ReadEventByCompanyId(c chan map[string]interface{}, processId string)
	ApplyBuildSteps()
	ApplyIntermediarySteps()
	ApplyJenkinsJobSteps()
	ApplyBuildCancellationSteps()
	PostOperationsForBuildPack(step string, stepType enums.STEP_TYPE, pipeline v1.Pipeline, claim int)
	CheckIfStepIsClaimable(step, processId string) bool
}
