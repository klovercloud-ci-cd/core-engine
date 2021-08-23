package v1

import (
	"github.com/tektoncd/pipeline/pkg/client/clientset/versioned"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1alpha1"
)
type Tekton interface {
	initPipelineResources(pipeline Pipeline) (input v1alpha1.PipelineResource,output v1alpha1.PipelineResource,err error)
	initTask(pipeline Pipeline) (v1alpha1.Task,error)
	initTaskRun(pipeline Pipeline)(v1alpha1.TaskRun,error)
}

type TektonResource struct {
	tcs *versioned.Clientset
}