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
	for _,each:=range p.pipeline.Steps{
		input,outputs,err:=p.tekton.InitPipelineResources(each)
		if err!=nil{return err}
		task,err:=p.tekton.InitTask(each)
		if err!=nil{return err}
		taskrun,err:=p.tekton.InitTaskRun(each)
		if err!=nil{return err}
	}



}

func NewPipelineService(k8s v1.K8s,tekton v1.Tekton,pipeline v1.Pipeline) service.Pipeline {
	return &pipelineService{
		k8s: k8s,
		tekton: tekton,
		pipeline: pipeline,
	}
}