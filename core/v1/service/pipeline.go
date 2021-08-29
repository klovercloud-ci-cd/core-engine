package service

import (
	v1 "github.com/klovercloud-ci/core/v1"
	"github.com/klovercloud-ci/enums"
)

type Pipeline interface {
	Apply(url,revision string,pipeline v1.Pipeline)error
	LoadArgs(pipeline v1.Pipeline)
	LoadEnvs(pipeline v1.Pipeline)
	SetInputResource(url,revision string,pipeline v1.Pipeline)
	Build(url,revision string,pipeline v1.Pipeline)
	PostOperations(revision,step string,stepType enums.STEP_TYPE,pipeline v1.Pipeline)
}