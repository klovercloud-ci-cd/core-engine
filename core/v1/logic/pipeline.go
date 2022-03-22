package logic

import (
	"errors"
	"github.com/klovercloud-ci-cd/core-engine/config"
	v1 "github.com/klovercloud-ci-cd/core-engine/core/v1"
	"github.com/klovercloud-ci-cd/core-engine/core/v1/service"
	"github.com/klovercloud-ci-cd/core-engine/enums"
	"log"
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

//ApplyBuildCancellationSteps pulls cancellation events for build job and the process the cancellation.
func (p *pipelineService) ApplyBuildCancellationSteps() {
	steps := p.processLifeCycleEvent.PullBuildCancellingEvents()
	for _, each := range steps {
		p.tekton.PurgeByProcessId(each.ProcessId)
		processEventData := make(map[string]interface{})
		processEventData["step"] = each.Step
		if each.Pipeline == nil {
			each.Pipeline = &v1.Pipeline{
				ProcessId: each.ProcessId,
			}
		}
		subject := v1.Subject{Pipeline: *each.Pipeline, Step: each.Step, Log: "Cancelled Successfully", StepType: enums.BUILD}
		go p.notifyAll(subject)
	}
}

//ApplyJenkinsJobSteps pulls jenkins steps and then applies.
func (p *pipelineService) ApplyJenkinsJobSteps() {
	events := p.processLifeCycleEvent.PullJenkinsJobStepsEvents()
	for _, each := range events {
		p.pipeline = *each.Pipeline
		for i, step := range each.Pipeline.Steps {
			if each.Step == step.Name && each.StepType == step.Type {
				p.applySteps(each.Pipeline.Steps[i])
			}
		}
	}
}

//ApplyIntermediarySteps pulls intermediary steps and then applies.
func (p *pipelineService) ApplyIntermediarySteps() {
	events := p.processLifeCycleEvent.PullIntermediaryStepsEvents()
	for _, each := range events {
		p.pipeline = *each.Pipeline
		for i, step := range each.Pipeline.Steps {
			if each.Step == step.Name && each.StepType == step.Type {
				p.applySteps(each.Pipeline.Steps[i])
			}
		}
	}
}

//ApplyBuildSteps pulls build steps and then applies.
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

//ReadEventByProcessId reads live events from queue and then dequeues read messages. This is used optionally, only when local event store is used.
func (p *pipelineService) ReadEventByProcessId(c chan map[string]interface{}, processId string) {
	c <- p.processEventService.DequeueByProcessId(processId)
}

//GetLogsByProcessId get Logs by process id and log event options. This is used optionally, only when local event store is used.
func (p *pipelineService) GetLogsByProcessId(processId string, option v1.LogEventQueryOption) ([]string, int64) {
	return p.logEventService.GetByProcessId(processId, option)
}

func (p *pipelineService) PostOperationsForBuildPack(step string, stepType enums.STEP_TYPE, pipeline v1.Pipeline) {
	podList := p.k8s.WaitAndGetInitializedPods(config.CiNamespace, pipeline.ProcessId, step)
	if len(podList.Items) > 0 {
		pod := podList.Items[0]
		for index := range pod.Spec.Containers {
			p.k8s.FollowContainerLifeCycle(pod.Namespace, pod.Name, pod.Spec.Containers[index].Name, step, pipeline.ProcessId, stepType)
		}
	}
	pRun, pRunError := p.tekton.GetPipelineRun(step+"-"+pipeline.ProcessId, true)
	pRunStatus := ""
	if pRunError != nil {
		pRunStatus = pRunError.Error()
	} else {
		if pRun.IsDone() {
			pRunStatus = string(enums.SUCCESSFUL)
		} else if pRun.IsCancelled() {
			pRunStatus = string(enums.CANCELLED)
		} else {
			pRunStatus = string(enums.ERROR)
		}
	}
	processEventData := make(map[string]interface{})
	processEventData["step"] = step
	processEventData["status"] = pRunStatus
	processEventData["type"] = stepType
	listener := v1.Subject{Pipeline: pipeline, Step: step}
	listener.EventData = processEventData
	go p.notifyAll(listener)
	if pipeline.Option.Purging == enums.PIPELINE_PURGING_ENABLE {
		go p.tekton.PurgeByProcessId(p.pipeline.ProcessId)
	}
}

// PostOperations Wait until pod is created, watches pod lifecycle, sends events to all the observers. Purges resources if purging is enabled.
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

//LoadArgs reads data from arg, serializes string into map and set into argData of step
func (p *pipelineService) LoadArgs(pipeline v1.Pipeline) {
	p.pipeline = pipeline
	for i := range p.pipeline.Steps {
		s := stepService{p.pipeline.Steps[i]}
		s.SetArgs(p.k8s)
		p.pipeline.Steps[i] = s.step
	}
}

// LoadEnvs reads data from env, serializes string into map and set into envData of step
func (p *pipelineService) LoadEnvs(pipeline v1.Pipeline) {
	p.pipeline = pipeline
	for i := range p.pipeline.Steps {
		s := stepService{p.pipeline.Steps[i]}
		s.SetEnvs(p.k8s)
		p.pipeline.Steps[i] = s.step
	}
}

// SetInputResource sets input resources for build step
func (p *pipelineService) SetInputResource(url, revision string, pipeline v1.Pipeline) {
	p.pipeline = pipeline
	for i := range p.pipeline.Steps {
		s := stepService{p.pipeline.Steps[i]}
		s.SetInput(url, revision)
		p.pipeline.Steps[i] = s.step
	}
}

// Build reads data from arg, serializes string into map and set into argData of step, reads data from env, serializes string into map and set into envData of step, sets input resources for build step
func (p *pipelineService) Build(url, revision string, pipeline v1.Pipeline) {
	p.LoadArgs(pipeline)
	p.LoadEnvs(pipeline)
	p.SetInputResource(url, revision, pipeline)
}

//BuildProcessLifeCycleEvents Builds pipeline By triggering Build method, then validates pipeline, and then notifies the observers.
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

// buildProcessLifeCycleEvents initializes build events subject and then notifies observers.
func (p *pipelineService) buildProcessLifeCycleEvents() {
	if len(p.pipeline.Steps) > 0 {
		initialStep := p.pipeline.Steps[0]
		processEventData := make(map[string]interface{})
		processEventData["step"] = initialStep.Name
		listener := v1.Subject{Pipeline: p.pipeline, Step: initialStep.Name}
		processEventData["trigger"] = initialStep.Trigger
		processEventData["status"] = enums.NON_INITIALIZED
		processEventData["next"] = strings.Join(initialStep.Next, ",")
		processEventData["type"] = initialStep.Type
		listener.EventData = processEventData
		if initialStep.Type == enums.BUILD || initialStep.Type == enums.INTERMEDIARY || initialStep.Type == enums.JENKINS_JOB {
			go p.notifyAll(listener)
		}
	}
}

// applySteps applies steps and then notifies observers.
func (p *pipelineService) applySteps(step v1.Step) {
	listener := v1.Subject{Pipeline: p.pipeline, Step: step.Name}
	processEventData := make(map[string]interface{})
	processEventData["step"] = step.Name
	processEventData["trigger"] = step.Params["trigger"]
	processEventData["type"] = step.Type
	var err error
	if step.Type == enums.BUILD {
		err = p.applyBuildStep(step)
	} else if step.Type == enums.INTERMEDIARY {
		err = p.applyIntermediaryStep(step)
	} else if step.Type == enums.JENKINS_JOB {
		err = p.applyJenkinsJobStep(step)
	} else {
		return
	}
	if err != nil {
		log.Println(err.Error())
		processEventData["status"] = string(enums.STEP_FAILED)
		processEventData["log"] = err
		listener.EventData = processEventData
		go p.notifyAll(listener)
		return
	}
	processEventData["status"] = string(enums.INITIALIZING)
	processEventData["next"] = strings.Join(step.Next, ",")
	listener.EventData = processEventData
	go p.notifyAll(listener)
}

//applyJenkinsJobStep applies jenkins step, follows pod lifecycle and the notifies observers
func (p *pipelineService) applyJenkinsJobStep(step v1.Step) error {
	trimmedStepName := strings.ReplaceAll(step.Name, " ", "")
	step.Name = trimmedStepName
	taskrun, err := p.tekton.InitTaskRun(step, p.pipeline.Label, p.pipeline.ProcessId)
	if err != nil {
		return errors.New("Failed to initialize pipeline job" + err.Error())
	}
	err = p.tekton.CreateTaskRun(taskrun)
	if err != nil {
		_ = p.tekton.DeleteTaskRunByProcessId(p.pipeline.ProcessId)
		return errors.New("Failed to apply taskrun" + err.Error())
	}
	go p.PostOperations(step.Name, step.Type, p.pipeline)
	return nil
}

//applyIntermediaryStep applies intermediary step, follows pod lifecycle and the notifies observers
func (p *pipelineService) applyIntermediaryStep(step v1.Step) error {
	trimmedStepName := strings.ReplaceAll(step.Name, " ", "")
	step.Name = trimmedStepName
	task, err := p.tekton.InitTask(step, p.pipeline.Label, p.pipeline.ProcessId)
	if err != nil {
		return errors.New("Failed to initialize task" + err.Error())
	}
	taskrun, err := p.tekton.InitTaskRun(step, p.pipeline.Label, p.pipeline.ProcessId)
	if err != nil {
		return errors.New("Failed to initialize pipeline job" + err.Error())
	}
	_ = p.tekton.DeleteTaskRunByProcessId(p.pipeline.ProcessId)
	if err != nil {
		return errors.New("Failed to apply input resource" + err.Error())
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

func (p *pipelineService) applyBuildPackStep(step v1.Step) error {
	// init pvc
	// init pipeline
	// init pipelinerun
	// then create pvc
	// create pipeline
	// create pipelinerun
	// if any step fails, delete all the resources by label
	pvc := p.k8s.InitPersistentVolumeClaim(step, p.pipeline.Label, p.pipeline.ProcessId)
	err := p.k8s.CreatePersistentVolumeClaim(pvc)
	if err != nil {
		return err
	}
	_ = p.k8s.DeletePersistentVolumeClaimByProcessId(p.pipeline.ProcessId)
	pipeline := p.tekton.InitPipeline(step, p.pipeline.Label, p.pipeline.ProcessId)
	err = p.tekton.CreatePipeline(pipeline)
	if err != nil {
		_ = p.tekton.DeletePipelineByProcessId(p.pipeline.ProcessId)
		return errors.New("Failed to apply pipeline" + err.Error())
	}
	pRun, err := p.tekton.InitPipelineRun(step, p.pipeline.Label, p.pipeline.ProcessId)
	if err != nil {
		return err
	}
	err = p.tekton.CreatePipelineRun(pRun)
	if err != nil {
		_ = p.tekton.DeletePipelineRunByProcessId(p.pipeline.ProcessId)
		return errors.New("Failed to apply pipeline run" + err.Error())
	}
	go p.PostOperations(step.Name, step.Type, p.pipeline)
	return nil
}

func (p *pipelineService) applyRegularBuildStep(step v1.Step) error {
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

//applyBuildStep applies build step, follows pod lifecycle and the notifies observers
func (p *pipelineService) applyBuildStep(step v1.Step) error {
	trimmedStepName := strings.ReplaceAll(step.Name, " ", "")
	step.Name = trimmedStepName
	// check if built type equals build back
	if step.Params[enums.BUILD_TYPE] == "buildpack" {
		return p.applyBuildPackStep(step)
	} else {
		return p.applyRegularBuildStep(step)
	}
}

// notifyAll notifies all the observers
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
