package logic

import (
	v1 "github.com/klovercloud-ci/core/v1"
	"github.com/klovercloud-ci/core/v1/service"
	"log"
)

type pipelineService struct {
	k8s v1.K8s
	tekton v1.Tekton
	pipeline v1.Pipeline
}

func (p pipelineService) apply() {
	// all the err logs will be persisted by buildId
	for _,each:=range p.pipeline.Steps{
		input,outputs,err:=p.tekton.InitPipelineResources(each,p.pipeline.Name,p.pipeline.Label,p.pipeline.BuildId)
		if err!=nil{
			log.Println(err.Error())
			continue
		}
		task,err:=p.tekton.InitTask(each,p.pipeline.Name,p.pipeline.Label,p.pipeline.BuildId)
		if err!=nil{
			log.Println(err.Error())
			continue
		}
		taskrun,err:=p.tekton.InitTaskRun(each,p.pipeline.Label, p.pipeline.BuildId)
		if err!=nil{
			log.Println(err.Error())
			continue
		}
		p.tekton.DeleteTaskRunByBuildId(p.pipeline.BuildId)
		err=p.tekton.CreatePipelineResource(input)
		if err!=nil{
			log.Println(err.Error())
			continue
		}
		var outputErr error
		for _,output:=range outputs{
			err=p.tekton.CreatePipelineResource(output)
			if err!=nil{
				outputErr=err
				break
			}
		}
		if outputErr!=nil{
			log.Println(outputErr.Error())
			p.tekton.DeleteTaskRunByBuildId(p.pipeline.BuildId)
			continue
		}
		err=p.tekton.CreateTask(task)
		if err!=nil{
			log.Println(err.Error())
			p.tekton.DeleteTaskRunByBuildId(p.pipeline.BuildId)
			continue
		}
		err=p.tekton.CreateTaskRun(taskrun)
		if err!=nil{
			log.Println(err.Error())
			p.tekton.DeleteTaskRunByBuildId(p.pipeline.BuildId)
			continue
		}

	}
}

func (p pipelineService) Apply(url,revision string) error {
	p.pipeline.Build(p.k8s,url,revision)
	//validate
	p.apply()
	return nil
}



func NewPipelineService(k8s v1.K8s,tekton v1.Tekton,pipeline v1.Pipeline) service.Pipeline {
	return &pipelineService{
		k8s: k8s,
		tekton: tekton,
		pipeline: pipeline,
	}
}