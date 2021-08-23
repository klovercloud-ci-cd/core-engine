package v1

import (
	"github.com/klovercloud-ci/config"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1alpha1"
	"github.com/tektoncd/pipeline/pkg/client/clientset/versioned"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)
type Tekton interface {
	InitPipelineResources(step Step) (input v1alpha1.PipelineResource,output []v1alpha1.PipelineResource,err error)
	InitTask(step Step) (v1alpha1.Task,error)
	InitTaskRun(step Step)(v1alpha1.TaskRun,error)
}

type TektonResource struct {
	Tcs  *versioned.Clientset

}

func (tekton * TektonResource) InitPipelineResources(step Step)(inputResource v1alpha1.PipelineResource,outputResource []v1alpha1.PipelineResource,err error){
	input:=v1alpha1.PipelineResource{
		TypeMeta:   v1.TypeMeta{},
		ObjectMeta: v1.ObjectMeta{
			Name:                       step.Name+"-"+step.Input.Revision,
			Namespace:                config.CI_NAMESPACE,
		},
	}
	if step.Input.Type==GIT{
		input.Spec.Type=string(GIT)
		input.Spec.Params=append(input.Spec.Params, struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		}{Name: "revision", Value:step.Input.Revision })
		input.Spec.Params=append(input.Spec.Params, struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		}{Name: "url", Value:step.Input.Url })
	}
	outputs:=[]v1alpha1.PipelineResource{}
	for i,_:=range step.Outputs{
		output:=v1alpha1.PipelineResource{
			TypeMeta:   v1.TypeMeta{},
			ObjectMeta: v1.ObjectMeta{
				Name:                       step.Name+"-"+step.Outputs[i].Revision,
				Namespace:                config.CI_NAMESPACE,
			},
		}
		if step.Outputs[i].Type==IMAGE{
			output.Spec.Type=string(IMAGE)
			output.Spec.Params=append(input.Spec.Params, struct {
				Name  string `json:"name"`
				Value string `json:"value"`
			}{Name: "url", Value:step.Outputs[i].Url })
		}
		outputs= append(outputs, output)
	}
	return input,outputs,nil
}

func (tekton * TektonResource) InitTask(step Step) (v1alpha1.Task,error){

	return v1alpha1.Task{},nil
}

func (tekton * TektonResource) InitTaskRun (step Step)(v1alpha1.TaskRun,error){

	return v1alpha1.TaskRun{},nil
}

