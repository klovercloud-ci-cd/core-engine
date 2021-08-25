package logic

import (
	"github.com/klovercloud-ci/config"
	"github.com/klovercloud-ci/core/v1/repository"
	"github.com/klovercloud-ci/core/v1/service"
	"github.com/klovercloud-ci/enums"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1alpha1"
	"github.com/tektoncd/pipeline/pkg/client/clientset/versioned"
	corev1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/klovercloud-ci/core/v1"
	"log"
	"strconv"
)
type tektonService struct {
	Tcs  *versioned.Clientset
	LogEventRepo repository.LogEventRepository
}

func (tekton *tektonService) InitPipelineResources(step v1.Step,label map[string]string,buildId string)(inputResource v1alpha1.PipelineResource,outputResource []v1alpha1.PipelineResource,err error){
	if label==nil{
		label=make(map[string]string)
	}
	label["revision"]=step.Input.Revision
	label["buildId"]=buildId
	input:=v1alpha1.PipelineResource{
		TypeMeta: metaV1.TypeMeta{
			Kind:       "PipelineResource",
			APIVersion: "tekton.dev/v1alpha1",
		},
		ObjectMeta: metaV1.ObjectMeta{
			Name:      step.Input.Revision + "-" + buildId,
			Namespace: config.CI_NAMESPACE,
			Labels:    label,
		},
	}
	if step.Input.Type==enums.GIT{
		input.Spec.Type= v1alpha1.PipelineResourceType(enums.GIT)
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
		label["step"]=step.Name
		label["revision"]=step.Outputs[i].Revision
		output:=v1alpha1.PipelineResource{
			ObjectMeta: metaV1.ObjectMeta{
				Name:                     step.Outputs[i].Revision+"-"+buildId,
				Namespace:                config.CI_NAMESPACE,
				Labels:    label,
			},
		}
		if step.Outputs[i].Type==enums.IMAGE{
			output.Spec.Type= v1alpha1.PipelineResourceType(enums.IMAGE)
			output.Spec.Params=append(input.Spec.Params, struct {
				Name  string `json:"name"`
				Value string `json:"value"`
			}{Name: "url", Value:step.Outputs[i].Url })
		}
		outputs= append(outputs, output)
	}
	return input,outputs,nil
}
func (tekton *tektonService) InitTask(step v1.Step,label map[string]string,buildId string) (v1alpha1.Task,error){
	if label==nil{
		label=make(map[string]string)
	}
	label["step"]=step.Name
	label["buildId"]=buildId
	task:=&v1alpha1.Task{
		TypeMeta: metaV1.TypeMeta{
			Kind:       "Task",
			APIVersion: "tekton.dev/v1",
		},
		ObjectMeta: metaV1.ObjectMeta{
			Name:      step.Name +"-"+buildId,
			Namespace: config.CI_NAMESPACE,
			Labels:    label,
		},
	}

	if step.Type==enums.BUILD{
		initBuildTaskSpec(step, task)
	}else {
		log.Print("Please provide a valid step type!")
	}

	return *task, nil

}
func initBuildTaskSpec(step v1.Step, task *v1alpha1.Task) {
	params := []v1alpha1.ParamSpec{}
	params = append(params, v1alpha1.ParamSpec{
		Name:        "pathToDockerFile",
		Type:        "string",
		Description: "The path to the dockerfile to build",
		Default: &v1alpha1.ArrayOrString{
			Type:      "string",
			StringVal: "/workspace/docker-source/Dockerfile",
		},
	})
	params = append(params, v1alpha1.ParamSpec{
		Name:        "pathToContext",
		Type:        "string",
		Description: "The build context used by Kaniko (https://github.com/GoogleContainerTools/kaniko#kaniko-build-contexts)",
		Default: &v1alpha1.ArrayOrString{
			Type:      "string",
			StringVal: "/workspace/docker-source",
		},
	})
	params = append(params, initAdditionalParams(step.Arg.Data)...)

	task.Spec.Inputs = &v1alpha1.Inputs{
		Resources: []v1alpha1.TaskResource{
			{ResourceDeclaration: v1alpha1.ResourceDeclaration{
				Name: "docker-source",
				Type: "git",
			}},
		},
		Params: params,
	}
	taskResource := []v1alpha1.TaskResource{}
	for i, _ := range step.Outputs {
		declaration := v1alpha1.ResourceDeclaration{
			Name: "builtImage" + strconv.Itoa(i),
			Type: "image",
		}
		resource := v1alpha1.TaskResource{declaration}
		taskResource = append(taskResource, resource)
	}
	task.Spec.Outputs = &v1alpha1.Outputs{
		Resources: taskResource,
	}
	var steps []v1alpha1.Step
	args := initBuildArgs(step.Arg.Data)
	args = append(args, "--dockerfile=$(inputs.params.pathToDockerFile)")
	args = append(args, "---context=$(inputs.params.pathToContext)")
	for i, _ := range step.Outputs {
		args = append(args, "--destination=$(outputs.resources.builtImage"+strconv.Itoa(i)+".url)")
		steps = append(steps, v1alpha1.Step{
			Container: corev1.Container{
				Name:    "build-and-push",
				Image:   enums.KLOVERCLOUD_KANIKO,
				Command: []string{"/kaniko/executor"},
				Args:    args,
				Env: []corev1.EnvVar{corev1.EnvVar{
					Name:  "DOCKER_CONFIG",
					Value: "/tekton/home/.docker/",
				}},
				ImagePullPolicy: "Always",
			},
		})
	}
	task.Spec.Steps = steps
}
func initBuildArgs(arg map[string]string) [] string{
	var args [] string
	for key,_:=range arg{
		args= append(args, "--build-arg="+key+"=$(inputs.params."+key+")")
	}
	return args
}
func initAdditionalParams(args map[string]string) []v1alpha1.ParamSpec {
	var params []v1alpha1.ParamSpec
	for key,_:=range args{
		params = append(params, v1alpha1.ParamSpec{
			Name:    key,
			Type:    "string",
			Default: nil,
		})
	}
	return params
}
func (tekton *tektonService) InitTaskRun (step v1.Step,label map[string]string,buildId string)(v1alpha1.TaskRun,error){
	label["step"]=step.Name
	var params []v1alpha1.Param
	params = append(params, v1alpha1.Param{
		Name:  "pathToDockerFile",
		Value: v1alpha1.ArrayOrString{
			Type:      v1alpha1.ParamTypeString,
			StringVal: "Dockerfile",
		},
	})
	params = append(params, v1alpha1.Param{
		Name:  "pathToContext",
		Value: v1alpha1.ArrayOrString{
			Type:      v1alpha1.ParamTypeString,
			StringVal: "/workspace/docker-source",
		},
	})

	taskrun:=v1alpha1.TaskRun{
		TypeMeta:   metaV1.TypeMeta{
			Kind:       "TaskRun",
			APIVersion: "tekton.dev/v1alpha1",
		},
		ObjectMeta: metaV1.ObjectMeta{
			Name:      step.Name +"-"+buildId,
			Namespace: config.CI_NAMESPACE,
			Labels:    label,
		},
	}

	if step.Type==enums.BUILD{
		if step.Arg.Data!=nil{
			for k,v:=range step.Arg.Data{
				params= append(params,v1alpha1.Param{
					Name:  k,
					Value: v1alpha1.ArrayOrString{
						Type:      v1alpha1.ParamTypeString,
						StringVal: v,
					},
				} )
			}

		}
		taskrun.Spec=v1alpha1.TaskRunSpec{
			ServiceAccountName: step.ServiceAccount,
			TaskRef:            &v1alpha1.TaskRef{
				Name: step.Name +"-"+buildId,
			},
		}

		taskrun.Spec.Inputs=v1alpha1.TaskRunInputs{
			Resources: []v1alpha1.TaskResourceBinding{v1alpha1.TaskResourceBinding{
				PipelineResourceBinding: v1alpha1.PipelineResourceBinding{
					Name: "docker-source",
					ResourceRef: &v1alpha1.PipelineResourceRef{
						Name: step.Input.Revision + "-" + buildId,
					},
				},
			}},
			Params:    params,
		}
		resourceBindings:=[]v1alpha1.TaskResourceBinding{}
		for i,each:=range step.Outputs{
			resourceBindings= append(resourceBindings, 	v1alpha1.TaskResourceBinding{
				PipelineResourceBinding: v1alpha1.PipelineResourceBinding{
					Name:         "builtImage" + strconv.Itoa(i),
					ResourceRef:  &v1alpha1.PipelineResourceRef{
						Name: each.Revision+"-"+buildId,
					},
				},
			})
		}
		taskRunOutputs:=v1alpha1.TaskRunOutputs{
			Resources: resourceBindings,
		}
		taskrun.Spec.Outputs=taskRunOutputs
	}

	return taskrun,nil
	return v1alpha1.TaskRun{},nil
}
func (tekton *tektonService) CreatePipelineResource(resource v1alpha1.PipelineResource)error {
	_,err:=tekton.Tcs.TektonV1alpha1().PipelineResources(config.CI_NAMESPACE).Create(&resource)
	if err!=nil{
		log.Println("[ERROR]:","Failed to create pipelineresource" ,err.Error())
		return err
	}
	return nil
}
func(tekton *tektonService)CreateTask(resource v1alpha1.Task)error{
	_,err:=tekton.Tcs.TektonV1alpha1().Tasks(config.CI_NAMESPACE).Create(&resource)
	if err!=nil{
		log.Println("[ERROR]:","Failed to create task" ,err.Error())
		return err
	}
	return nil
}
func(tekton *tektonService)CreateTaskRun(resource v1alpha1.TaskRun) error{
	_,err:=tekton.Tcs.TektonV1alpha1().TaskRuns(config.CI_NAMESPACE).Create(&resource)
	if err!=nil{
		log.Println("[ERROR]:","Failed to create taskrun" ,err.Error())
		return err
	}
	return nil

}
func(tekton *tektonService)DeletePipelineResourceByBuildId(buildId string) error{
	list,err:=tekton.Tcs.TektonV1alpha1().PipelineResources(config.CI_NAMESPACE).List(metaV1.ListOptions{
		LabelSelector: "buildId="+buildId,
	})
	if err!=nil{
		log.Println("[WARNING]:",err.Error())
		return err
	}
	for _,each:=range list.Items{
		err=tekton.Tcs.TektonV1alpha1().PipelineResources(config.CI_NAMESPACE).Delete(each.Name,&metaV1.DeleteOptions{})
		if err!=nil{
			log.Println("[ERROR]:",err.Error())
		}
	}
	return nil
}
func(tekton *tektonService)DeleteTaskByBuildId(buildId string) error{
	list,err:=tekton.Tcs.TektonV1alpha1().Tasks(config.CI_NAMESPACE).List(metaV1.ListOptions{
		LabelSelector: "buildId="+buildId,
	})
	if err!=nil{
		log.Println("[WARNING]:",err.Error())
		return err
	}
	for _,each:=range list.Items{
		err=tekton.Tcs.TektonV1alpha1().Tasks(config.CI_NAMESPACE).Delete(each.Name,&metaV1.DeleteOptions{})
		if err!=nil{
			log.Println("[ERROR]:",err.Error())
		}
	}
	return nil
}
func(tekton *tektonService)DeleteTaskRunByBuildId(buildId string) error{
	list,err:=tekton.Tcs.TektonV1alpha1().TaskRuns(config.CI_NAMESPACE).List(metaV1.ListOptions{
		LabelSelector: "buildId="+buildId,
	})
	if err!=nil{
		log.Println("[WARNING]:",err.Error())
		return err
	}
	for _,each:=range list.Items{
		err=tekton.Tcs.TektonV1alpha1().TaskRuns(config.CI_NAMESPACE).Delete(each.Name,&metaV1.DeleteOptions{})
		if err!=nil{
			log.Println("[ERROR]:",err.Error())
		}
	}
	return nil
}
func(tekton *tektonService)PurgeByBuildId(buildId string) {
	tekton.DeletePipelineResourceByBuildId(buildId)
	tekton.DeleteTaskByBuildId(buildId)
	tekton.DeleteTaskRunByBuildId(buildId)
}

func NewTektonService(tcs  *versioned.Clientset,logEventRepo repository.LogEventRepository) service.Tekton{
	return  &tektonService{
		Tcs:          tcs,
		LogEventRepo:logEventRepo,
	}
}