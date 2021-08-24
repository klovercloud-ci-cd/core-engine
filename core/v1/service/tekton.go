package service

import (
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1alpha1"
	"github.com/klovercloud-ci/core/v1"
)

type Tekton interface {
	InitPipelineResources(step v1.Step,label map[string]string,buildId string) (input v1alpha1.PipelineResource,output []v1alpha1.PipelineResource,err error)
	InitTask(step v1.Step,label map[string]string,buildId string) (v1alpha1.Task,error)
	InitTaskRun(step v1.Step,label map[string]string,buildId string)(v1alpha1.TaskRun,error)
	CreatePipelineResource(v1alpha1.PipelineResource) error
	CreateTask(v1alpha1.Task) error
	CreateTaskRun(v1alpha1.TaskRun) error
	DeletePipelineResourceByBuildId(buildId string) error
	DeleteTaskByBuildId(buildId string) error
	DeleteTaskRunByBuildId(buildId string) error
	PurgeByBuildId(buildId string)
}