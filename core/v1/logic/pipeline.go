package logic

import (
	v1 "github.com/klovercloud-ci/core/v1"
	"github.com/klovercloud-ci/core/v1/service"
)

type pipelineService struct {
	k8s v1.K8s
	tekton v1.Tekton
	pipeline v1.Pipeline
}

func (p pipelineService) Apply(url,revision string) error {
	p.pipeline.Build(p.k8s,url,revision)
	panic("implement me")
}

func NewPipelineService(k8s v1.K8s) service.Pipeline {
	return &pipelineService{
		k8s: k8s,
	}
}