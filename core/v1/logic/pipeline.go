package logic

import (
	"errors"
	"github.com/klovercloud-ci-cd/core-engine/config"
	v1 "github.com/klovercloud-ci-cd/core-engine/core/v1"
	"github.com/klovercloud-ci-cd/core-engine/core/v1/service"
	"github.com/klovercloud-ci-cd/core-engine/enums"
	"strings"
)

type pipelineService struct {
	k8s                   service.K8s
	tekton                service.Tekton
	pipeline              v1.Pipeline
	logEventService       service.LogEvent
	processEventService   service.ProcessEvent
	processLifeCycleEvent service.ProcessLifeCycleEvent
	observerList          []service.Observer
}

func (p *pipelineService) ApplyBuildSteps() {
	events := p.processLifeCycleEvent.PullBuildEvents()
	for _, each := range events {
		p.pipeline = *each.Pipeline
		for i, step := range each.Pipeline.Steps {
			if each.Step == step.Name && each.StepType == step.Type {
				p.applySteps(each.Pipeline.Steps[i])
			}
		}
	}
}

func (p *pipelineService) ReadEventByProcessId(c chan map[string]interface{}, processId string) {
	c <- p.processEventService.DequeueByProcessId(processId)
}

func (p *pipelineService) GetLogsByProcessId(processId string, option v1.LogEventQueryOption) ([]string, int64) {
	return p.logEventService.GetByProcessId(processId, option)
}

func (p *pipelineService) PostOperations(step string, stepType enums.STEP_TYPE, pipeline v1.Pipeline) {
	podList := p.k8s.WaitAndGetInitializedPods(config.CiNamespace, pipeline.ProcessId, step)
	if len(podList.Items) > 0 {
		pod := podList.Items[0]
		for index := range pod.Spec.Containers {
			p.k8s.FollowContainerLifeCycle(pod.Namespace, pod.Name, pod.Spec.Containers[index].Name, step, pipeline.ProcessId, stepType)
		}
	}
	tRun, tRunError := p.tekton.GetTaskRun(step+"-"+pipeline.ProcessId, true)
	tRunStatus := ""
	if tRunError != nil {
		tRunStatus = tRunError.Error()
	} else {
		if tRun.IsSuccessful() {
			tRunStatus = string(enums.SUCCESSFUL)
		} else if tRun.IsCancelled() {
			tRunStatus = string(enums.CANCELLED)
		} else {
			tRunStatus = string(enums.ERROR)
		}
	}
	processEventData := make(map[string]interface{})
	processEventData["step"] = step
	processEventData["status"] = tRunStatus
	processEventData["type"] = stepType
	listener := v1.Subject{Pipeline: pipeline, Step: step}
	listener.EventData = processEventData
	go p.notifyAll(listener)
	if pipeline.Option.Purging == enums.PIPELINE_PURGING_ENABLE {
		go p.tekton.PurgeByProcessId(p.pipeline.ProcessId)
	}
}

func (p *pipelineService) LoadArgs(pipeline v1.Pipeline) {
	p.pipeline = pipeline
	for i := range p.pipeline.Steps {
		s := stepService{p.pipeline.Steps[i]}
		s.SetArgs(p.k8s)
		p.pipeline.Steps[i] = s.step
	}
}

func (p *pipelineService) LoadEnvs(pipeline v1.Pipeline) {
	p.pipeline = pipeline
	for i := range p.pipeline.Steps {
		s := stepService{p.pipeline.Steps[i]}
		s.SetEnvs(p.k8s)
		p.pipeline.Steps[i] = s.step
	}
}

func (p *pipelineService) SetInputResource(url, revision string, pipeline v1.Pipeline) {
	p.pipeline = pipeline
	for i := range p.pipeline.Steps {
		s := stepService{p.pipeline.Steps[i]}
		s.SetInput(url, revision)
		p.pipeline.Steps[i] = s.step
	}
}

func (p *pipelineService) Build(url, revision string, pipeline v1.Pipeline) {
	p.LoadArgs(pipeline)
	p.LoadEnvs(pipeline)
	p.SetInputResource(url, revision, pipeline)
}

func (p *pipelineService) BuildProcessLifeCycleEvents(url, revision string, pipeline v1.Pipeline) error {
	p.Build(url, revision, pipeline)
	if p.pipeline.Label == nil {
		p.pipeline.Label = make(map[string]string)
	}
	p.pipeline.Label["processId"] = p.pipeline.ProcessId
	p.pipeline.Label["pipeline"] = p.pipeline.Name
	err := p.pipeline.Validate()
	if err != nil {
		return err
	}
	p.buildProcessLifeCycleEvents()
	return nil
}
func (p *pipelineService) buildProcessLifeCycleEvents() {
	if len(p.pipeline.Steps) > 0 {
		initialStep := p.pipeline.Steps[0]
		processEventData := make(map[string]interface{})
		processEventData["step"] = initialStep.Name
		listener := v1.Subject{Pipeline: p.pipeline, Step: initialStep.Name}
		if initialStep.Type == enums.BUILD {
			processEventData["trigger"] = initialStep.Trigger
			processEventData["agent"] = initialStep.Params[enums.AGENT]
			processEventData["type"] = enums.BUILD
			processEventData["status"] = enums.NON_INITIALIZED
			processEventData["next"] = strings.Join(initialStep.Next, ",")
			listener.EventData = processEventData
			go p.notifyAll(listener)
		}
	}
}
func (p *pipelineService) applySteps(step v1.Step) {

	processEventData := make(map[string]interface{})
	processEventData["step"] = step.Name
	listener := v1.Subject{Pipeline: p.pipeline, Step: step.Name}
	if step.Type == enums.BUILD {
		err := p.applyBuildStep(step)
		processEventData["trigger"] = step.Params["trigger"]
		processEventData["agent"] = step.Params[enums.AGENT]
		processEventData["type"] = enums.BUILD
		if err != nil {
			processEventData["status"] = enums.BUILD_FAILED
			processEventData["log"] = err
			listener.EventData = processEventData
			go p.notifyAll(listener)
			return
		}

		processEventData["status"] = enums.INITIALIZING
		processEventData["next"] = strings.Join(step.Next, ",")
		listener.EventData = processEventData
		go p.notifyAll(listener)
	}
}
func (p *pipelineService) apply() {
	if len(p.pipeline.Steps) > 0 {
		initialStep := p.pipeline.Steps[0]
		processEventData := make(map[string]interface{})
		processEventData["step"] = initialStep.Name
		listener := v1.Subject{Pipeline: p.pipeline, Step: initialStep.Name}
		if initialStep.Type == enums.BUILD {
			err := p.applyBuildStep(initialStep)
			processEventData["trigger"] = initialStep.Params["trigger"]
			processEventData["agent"] = initialStep.Params[enums.AGENT]
			processEventData["type"] = enums.BUILD
			if err != nil {
				processEventData["status"] = enums.BUILD_FAILED
				processEventData["log"] = err
				listener.EventData = processEventData
				go p.notifyAll(listener)
				return
			}

			processEventData["status"] = enums.INITIALIZING
			processEventData["next"] = strings.Join(initialStep.Next, ",")
			listener.EventData = processEventData
			go p.notifyAll(listener)
		}
	}
}

func (p *pipelineService) applyBuildStep(step v1.Step) error {
	nss := strings.ReplaceAll(step.Name, " ", "")
	step.Name = nss
	input, outputs, err := p.tekton.InitPipelineResources(step, p.pipeline.Label, p.pipeline.ProcessId)
	if err != nil {
		return errors.New("Failed to initialize input/output resource" + err.Error())
	}
	task, err := p.tekton.InitTask(step, p.pipeline.Label, p.pipeline.ProcessId)
	if err != nil {
		return errors.New("Failed to initialize task" + err.Error())
	}
	taskrun, err := p.tekton.InitTaskRun(step, p.pipeline.Label, p.pipeline.ProcessId)
	if err != nil {
		return errors.New("Failed to initialize pipeline job" + err.Error())
	}
	_ = p.tekton.DeleteTaskRunByProcessId(p.pipeline.ProcessId)
	err = p.tekton.CreatePipelineResource(input)
	if err != nil {
		return errors.New("Failed to apply input resource" + err.Error())
	}
	var outputErr error
	for _, output := range outputs {
		err = p.tekton.CreatePipelineResource(output)
		if err != nil {
			outputErr = err
			break
		}
	}
	if outputErr != nil {
		_ = p.tekton.DeleteTaskRunByProcessId(p.pipeline.ProcessId)
		return errors.New("Failed to apply output resource" + outputErr.Error())
	}
	err = p.tekton.CreateTask(task)
	if err != nil {
		_ = p.tekton.DeleteTaskRunByProcessId(p.pipeline.ProcessId)
		return errors.New("Failed to apply task" + err.Error())
	}
	err = p.tekton.CreateTaskRun(taskrun)
	if err != nil {
		_ = p.tekton.DeleteTaskRunByProcessId(p.pipeline.ProcessId)
		return errors.New("Failed to apply taskrun" + err.Error())
	}
	go p.PostOperations(step.Name, step.Type, p.pipeline)
	return nil
}

func (p *pipelineService) notifyAll(listener v1.Subject) {
	for _, observer := range p.observerList {
		go observer.Listen(listener)
	}
}

// NewPipelineService returns Pipeline type service
func NewPipelineService(k8s service.K8s, tekton service.Tekton, logEventService service.LogEvent, processEventService service.ProcessEvent, observerList []service.Observer, processLifeCycleEvent service.ProcessLifeCycleEvent) service.Pipeline {
	return &pipelineService{
		k8s:                   k8s,
		tekton:                tekton,
		pipeline:              v1.Pipeline{},
		logEventService:       logEventService,
		processEventService:   processEventService,
		observerList:          observerList,
		processLifeCycleEvent: processLifeCycleEvent,
	}
}
