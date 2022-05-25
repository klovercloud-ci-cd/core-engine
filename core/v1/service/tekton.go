package service

import (
	"github.com/klovercloud-ci-cd/core-engine/core/v1"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1alpha1"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	corev1 "k8s.io/api/core/v1"
)

// Tekton tekton related operations.
type Tekton interface {
	InitPipelineResources(step v1.Step, label map[string]string, processId string) (input v1alpha1.PipelineResource, output []v1alpha1.PipelineResource, err error)
	InitTask(step v1.Step, label map[string]string, processId string) (v1alpha1.Task, error)
	InitTaskRun(step v1.Step, label map[string]string, processId string) (v1alpha1.TaskRun, error)
	CreatePipelineResource(v1alpha1.PipelineResource) error
	CreateTask(v1alpha1.Task) error
	CreateTaskRun(v1alpha1.TaskRun) error
	DeletePipelineResourceByProcessId(processId string) error
	DeleteTaskByProcessId(processId string) error
	DeletePipelineByProcessId(processId string) error
	DeletePipelineRunByProcessId(processId string) error
	DeleteTaskRunByProcessId(processId string) error
	PurgeByProcessId(processId string)
	GetTaskRun(name string, waitUntilTaskRunIsCompleted bool) (*v1alpha1.TaskRun, error)
	GetPipelineRun(companyId,name, id, stepType string, waitUntilPipelineRunIsCompleted bool, podList corev1.PodList, claim int) (*v1beta1.PipelineRun, error)
	CreatePipeline(pipeline v1beta1.Pipeline) error
	InitPipeline(step v1.Step, label map[string]string, processId string) v1beta1.Pipeline
	CreatePipelineRun(pipelineRun v1beta1.PipelineRun) error
	InitPipelineRun(step v1.Step, label map[string]string, processId string) (v1beta1.PipelineRun, error)
}
