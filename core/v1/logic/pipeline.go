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

func (p *pipelineService) LoadArgs(pipeline v1.Pipeline) {
	p.pipeline=pipeline
	for i,_:=range p.pipeline.Steps {
		s := stepService{p.pipeline.Steps[i]}
		s.SetArgs(p.k8s)
		p.pipeline.Steps[i]=s.step
	}

}

func (p *pipelineService) LoadEnvs(pipeline v1.Pipeline) {
	p.pipeline=pipeline
	for i,_:=range p.pipeline.Steps{
		s := stepService{p.pipeline.Steps[i]}
		s.SetEnvs(p.k8s)
		p.pipeline.Steps[i]=s.step
	}
}

func (p *pipelineService) SetInputResource(url, revision string,pipeline v1.Pipeline) {
	p.pipeline=pipeline
	for i,_:=range p.pipeline.Steps{
		s := stepService{p.pipeline.Steps[i]}
		s.SetInput(url,revision)
		p.pipeline.Steps[i]=s.step
	}
}

func (p *pipelineService) Build( url, revision string,pipeline v1.Pipeline) {
	p.LoadArgs(pipeline)
	p.LoadEnvs(pipeline)
	p.SetInputResource(url,revision,pipeline)
}

func (p *pipelineService) apply() {
	// all the err logs will be persisted by processId
	for _,each:=range p.pipeline.Steps{
		input,outputs,err:=p.tekton.InitPipelineResources(each,p.pipeline.Label,p.pipeline.ProcessId)
		if err!=nil{
			log.Println(err.Error())
			continue
		}
		task,err:=p.tekton.InitTask(each,p.pipeline.Label,p.pipeline.ProcessId)
		if err!=nil{
			log.Println(err.Error())
			continue
		}
		taskrun,err:=p.tekton.InitTaskRun(each,p.pipeline.Label, p.pipeline.ProcessId)
		if err!=nil{
			log.Println(err.Error())
			continue
		}
		p.tekton.DeleteTaskRunByProcessId(p.pipeline.ProcessId)
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
			p.tekton.DeleteTaskRunByProcessId(p.pipeline.ProcessId)
			continue
		}
		err=p.tekton.CreateTask(task)
		if err!=nil{
			log.Println(err.Error())
			p.tekton.DeleteTaskRunByProcessId(p.pipeline.ProcessId)
			continue
		}
		err=p.tekton.CreateTaskRun(taskrun)
		if err!=nil{
			log.Println(err.Error())
			p.tekton.DeleteTaskRunByProcessId(p.pipeline.ProcessId)
			continue
		}

	}
}

func (p *pipelineService) Apply(url,revision string,pipeline v1.Pipeline) error {
	p.Build(url,revision,pipeline)
	if p.pipeline.Label==nil{
		p.pipeline.Label=make(map[string]string)
	}
	p.pipeline.Label["processId"]=p.pipeline.ProcessId
	p.pipeline.Label["pipeline"]=p.pipeline.Name
	err:=p.pipeline.Validate()
	if err!=nil{
		return err
	}
	p.apply()
	return nil
}


func NewPipelineService(k8s service.K8s,tekton service.Tekton) service.Pipeline {
	return &pipelineService{
		k8s: k8s,
		tekton: tekton,
		pipeline: v1.Pipeline{},
	}
}