package logic

import (
	v1 "github.com/klovercloud-ci/core/v1"
	"github.com/klovercloud-ci/core/v1/service"
	"log"
)

type pipelineService struct {
	k8s service.K8s
	tekton service.Tekton
	pipeline v1.Pipeline
}

func (p pipelineService) LoadArgs(k8s service.K8s) {
	for i,_:=range p.pipeline.Steps{
		s := stepService{p.pipeline.Steps[i]}
		s.SetArgs(k8s)
		p.pipeline.Steps[i]=s.step
	}

}

func (p pipelineService) LoadEnvs(k8s service.K8s) {
	for i,_:=range p.pipeline.Steps{
		s := stepService{p.pipeline.Steps[i]}
		s.SetEnvs(k8s)
		p.pipeline.Steps[i]=s.step
	}
}

func (p pipelineService) SetInputResource(url, revision string) {
	for i,_:=range p.pipeline.Steps{
		s := stepService{p.pipeline.Steps[i]}
		s.SetInput(url,revision)
		p.pipeline.Steps[i]=s.step
	}
}

func (p pipelineService) Build(k8s service.K8s, url, revision string) {
	p.LoadArgs(k8s)
	p.LoadEnvs(k8s)
	p.SetInputResource(url,revision)
}

func (p pipelineService) apply() {
	// all the err logs will be persisted by buildId
	for _,each:=range p.pipeline.Steps{
		input,outputs,err:=p.tekton.InitPipelineResources(each,p.pipeline.Label,p.pipeline.BuildId)
		if err!=nil{
			log.Println(err.Error())
			continue
		}
		task,err:=p.tekton.InitTask(each,p.pipeline.Label,p.pipeline.BuildId)
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
	p.Build(p.k8s,url,revision)
	if p.pipeline.Label==nil{
		p.pipeline.Label=make(map[string]string)
	}
	p.pipeline.Label["buildId"]=p.pipeline.BuildId
	p.pipeline.Label["pipeline"]=p.pipeline.Name
	err:=p.pipeline.Validate()
	if err!=nil{
		return err
	}
	p.apply()
	return nil
}


func NewPipelineService(k8s service.K8s,tekton service.Tekton,pipeline v1.Pipeline) service.Pipeline {
	return &pipelineService{
		k8s: k8s,
		tekton: tekton,
		pipeline: pipeline,
	}
}