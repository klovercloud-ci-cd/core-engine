package logic

import (
	"errors"
	"github.com/klovercloud-ci/config"
	v1 "github.com/klovercloud-ci/core/v1"
	"github.com/klovercloud-ci/core/v1/service"
	"github.com/klovercloud-ci/enums"
	"strings"
)

type pipelineService struct {
	k8s service.K8s
	tekton service.Tekton
	pipeline v1.Pipeline
	logEventService service.LogEvent
	processEventService service.ProcessEvent
	observerList []service.Observer
}

func (p *pipelineService) ReadEventByProcessId(c chan map[string]interface{},processId string)  {
	c<- p.processEventService.DequeueByProcessId(processId)
}

func (p *pipelineService) GetLogsByProcessId(processId string, option v1.LogEventQueryOption) ([]string, int64) {
	return p.logEventService.GetByProcessId(processId,option)
}

func (p *pipelineService) PostOperations(revision,step  string,stepType enums.STEP_TYPE, pipeline v1.Pipeline) {
	podList:=p.k8s.WaitAndGetInitializedPods(config.CiNamespace,pipeline.ProcessId,step)
	if len(podList.Items) > 0 {
		pod := podList.Items[0]
		for index := range pod.Spec.Containers {
			 p.k8s.FollowContainerLifeCycle(pod.Namespace,pod.Name, pod.Spec.Containers[index].Name,step,pipeline.ProcessId,stepType)
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
			tRunStatus=string(enums.ERROR)
		}
	}
	processEventData :=make(map[string]interface{})
	processEventData["step"]=step
	processEventData["status"]=tRunStatus
	listener:=v1.Listener{ProcessId: p.pipeline.ProcessId,Step: step}
	listener.EventData=processEventData
	go p.notifyAll(listener)

	if pipeline.Option.Purging==enums.PIPELINE_PURGING_ENABLE{
		go p.tekton.PurgeByProcessId(p.pipeline.ProcessId)
	}
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
		processEventData :=make(map[string]interface{})
		processEventData["step"]=each.Name
		listener:=v1.Listener{ProcessId: p.pipeline.ProcessId,Step: each.Name}
		if each.Type==enums.BUILD{
			err:=p.applyBuildStep(each)
			if err!=nil{
				processEventData["status"]=enums.BUILD_FAILED
				processEventData["log"]=err
				listener.EventData=processEventData
				go p.notifyAll(listener)
				break
			}

			processEventData["status"]=enums.INITIALIZING
			listener.EventData=processEventData
			go p.notifyAll(listener)
		}
	}
}

func (p *pipelineService) applyBuildStep(step v1.Step) error{
	nss := strings.ReplaceAll(step.Name, " ", "")
	step.Name = nss
	input,outputs,err:=p.tekton.InitPipelineResources(step,p.pipeline.Label,p.pipeline.ProcessId)
	if err!=nil{
		return errors.New("Failed to initialize input/output resource"+err.Error())
	}
	task,err:=p.tekton.InitTask(step,p.pipeline.Label,p.pipeline.ProcessId)
	if err!=nil{
		return errors.New("Failed to initialize task"+err.Error())
	}
	taskrun,err:=p.tekton.InitTaskRun(step,p.pipeline.Label, p.pipeline.ProcessId)
	if err!=nil{
		return errors.New("Failed to initialize pipeline job"+err.Error())
	}
	_ = p.tekton.DeleteTaskRunByProcessId(p.pipeline.ProcessId)
	err=p.tekton.CreatePipelineResource(input)
	if err!=nil{
		return errors.New("Failed to apply input resource"+err.Error())
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
		_ = p.tekton.DeleteTaskRunByProcessId(p.pipeline.ProcessId)
		return errors.New("Failed to apply output resource"+outputErr.Error())
	}
	err=p.tekton.CreateTask(task)
	if err!=nil{
		_ = p.tekton.DeleteTaskRunByProcessId(p.pipeline.ProcessId)
		return errors.New("Failed to apply task"+err.Error())
	}
	err=p.tekton.CreateTaskRun(taskrun)
	if err!=nil{
		_ = p.tekton.DeleteTaskRunByProcessId(p.pipeline.ProcessId)
		return errors.New("Failed to apply taskrun"+err.Error())
	}
	go p.PostOperations(step.Input.Revision,step.Name,step.Type,p.pipeline)
	return nil
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
func (k8s *pipelineService)notifyAll(listener v1.Listener){
	for _, observer := range k8s.observerList {
		go observer.Listen(listener)
	}
}

func NewPipelineService(k8s service.K8s,tekton service.Tekton,logEventService service.LogEvent,	processEventService service.ProcessEvent,observerList []service.Observer) service.Pipeline {
	return &pipelineService{
		k8s: k8s,
		tekton: tekton,
		pipeline: v1.Pipeline{},
		logEventService: logEventService,
		processEventService: processEventService,
		observerList: observerList,
	}
}