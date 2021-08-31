package logic

import (
	"bytes"
	"github.com/klovercloud-ci/config"
	v1 "github.com/klovercloud-ci/core/v1"
	"github.com/klovercloud-ci/core/v1/service"
	"github.com/klovercloud-ci/enums"
	"io"
	"log"
	"strings"
)

type pipelineService struct {
	k8s service.K8s
	tekton service.Tekton
	pipeline v1.Pipeline
	logEventService service.LogEvent
	processEventService service.ProcessEvent
}

func (p *pipelineService) ReadEventByProcessId(processId string) map[string]interface{} {
	return p.processEventService.DequeueByProcessId(processId)
}

func (p *pipelineService) GetLogsByProcessId(processId string, option v1.LogEventQueryOption) ([]string, int64) {
	return p.logEventService.GetByProcessId(processId,option)
}

func (p *pipelineService) PostOperations(revision,step  string,stepType enums.STEP_TYPE, pipeline v1.Pipeline) {
	var  buf bytes.Buffer
	podList:=p.k8s.WaitAndGetInitializedPods(config.CiNamespace,pipeline.ProcessId,step)
	if len(podList.Items) > 0 {
		pod := podList.Items[0]
		for index := range pod.Spec.Containers {
			go p.k8s.FollowContainerLifeCycle(pod.Namespace,pod.Name, pod.Spec.Containers[index].Name,step,pipeline.ProcessId,stepType)
		}
		for index := range pod.Spec.Containers {
			readCloser, _ := p.k8s.GetContainerLog(pod.Namespace,pod.Name,pod.Spec.Containers[index].Name,pod.Labels)
			if readCloser != nil {
				var tempBuf bytes.Buffer
				io.Copy(&tempBuf, readCloser)
				buf.WriteString(tempBuf.String())
				readCloser.Close()
			}
		}
	}
	tRun, tRunError :=p.tekton.GetTaskRun(step + "-" + pipeline.ProcessId,true)
	tRunStatus := ""
	if tRunError != nil {
		tRunStatus = tRunError.Error()
	}else{
		if tRun.IsSuccessful(){
			tRunStatus=string(enums.SUCCESSFUL)

		} else if tRun.IsCancelled(){
			tRunStatus=string(enums.CANCELLED)
		}else {

		}
	}
	log.Println(tRunStatus)
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
	for _,each:=range p.pipeline.Steps{
		nss := strings.ReplaceAll(each.Name, " ", "")
		each.Name = nss
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
		_ = p.tekton.DeleteTaskRunByProcessId(p.pipeline.ProcessId)
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
			_ = p.tekton.DeleteTaskRunByProcessId(p.pipeline.ProcessId)
			continue
		}
		err=p.tekton.CreateTask(task)
		if err!=nil{
			log.Println(err.Error())
			_ = p.tekton.DeleteTaskRunByProcessId(p.pipeline.ProcessId)
			continue
		}
		err=p.tekton.CreateTaskRun(taskrun)
		if err!=nil{
			log.Println(err.Error())
			_ = p.tekton.DeleteTaskRunByProcessId(p.pipeline.ProcessId)
			continue
		}
	  go p.PostOperations(each.Input.Revision,each.Name,each.Type,p.pipeline)
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
	go p.apply()
	return nil
}


func NewPipelineService(k8s service.K8s,tekton service.Tekton,logEventService service.LogEvent,	processEventService service.ProcessEvent) service.Pipeline {
	return &pipelineService{
		k8s: k8s,
		tekton: tekton,
		pipeline: v1.Pipeline{},
		logEventService: logEventService,
		processEventService: processEventService,
	}
}