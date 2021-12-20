package logic

import (
	"context"
	"github.com/klovercloud-ci-cd/core-engine/config"
	"github.com/klovercloud-ci-cd/core-engine/core/v1"
	"github.com/klovercloud-ci-cd/core-engine/core/v1/service"
	"github.com/klovercloud-ci-cd/core-engine/enums"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1alpha1"
	"github.com/tektoncd/pipeline/pkg/client/clientset/versioned"
	corev1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
	"strconv"
	"strings"
)

type tektonService struct {
	Tcs *versioned.Clientset
}

func (tekton *tektonService) GetTaskRun(name string, waitUntilTaskRunIsCompleted bool) (*v1alpha1.TaskRun, error) {
	tRun, taskrunGetingErr := tekton.Tcs.TektonV1alpha1().TaskRuns(config.CiNamespace).Get(name, metaV1.GetOptions{
		TypeMeta: metaV1.TypeMeta{
			Kind:       "TaskRun",
			APIVersion: "tekton.dev/v1",
		},
	})
	if taskrunGetingErr != nil {
		return nil, taskrunGetingErr
	}
	if !tRun.IsDone() && waitUntilTaskRunIsCompleted == true {
		return tekton.GetTaskRun(name, waitUntilTaskRunIsCompleted)
	}
	return tRun, nil
}

func (tekton *tektonService) InitPipelineResources(step v1.Step, label map[string]string, processId string) (inputResource v1alpha1.PipelineResource, outputResource []v1alpha1.PipelineResource, err error) {
	if label == nil {
		label = make(map[string]string)
	}
	label["revision"] = step.Params[enums.REVISION]
	label["processId"] = processId
	input := v1alpha1.PipelineResource{
		TypeMeta: metaV1.TypeMeta{
			Kind:       "PipelineResource",
			APIVersion: "tekton.dev/v1alpha1",
		},
		ObjectMeta: metaV1.ObjectMeta{
			Name:      step.Params[enums.REVISION][:15] + "-" + processId,
			Namespace: config.CiNamespace,
			Labels:    label,
		},
	}
	if step.Params[enums.REPOSITORY_TYPE] == string(enums.GIT) {
		input.Spec.Type = v1alpha1.PipelineResourceType(enums.GIT)
		input.Spec.Params = append(input.Spec.Params, struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		}{Name: "revision", Value: step.Params[enums.REVISION]})
		input.Spec.Params = append(input.Spec.Params, struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		}{Name: "url", Value: step.Params[enums.IMAGE_URL]})
	}
	error := input.Validate(context.Background())
	if error != nil {
		return input, nil, error
	}

	outputs := []v1alpha1.PipelineResource{}

	for i, each := range strings.Split(step.Params[enums.IMAGES], ",") {
		attrs := strings.Split(each, ":")
		imageRevision := "latest"
		if len(attrs) == 2 {
			imageRevision = attrs[1]
		}
		label["revision"] = imageRevision
		label["step"] = step.Name
		output := v1alpha1.PipelineResource{
			ObjectMeta: metaV1.ObjectMeta{
				Name:      step.Name + "" + processId + "" + strconv.Itoa(i),
				Namespace: config.CiNamespace,
				Labels:    label,
			},
		}
		output.Spec.Type = v1alpha1.PipelineResourceType(enums.IMAGE)
		output.Spec.Params = append(output.Spec.Params, struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		}{Name: "url", Value: each})
		output.Spec.Params = append(output.Spec.Params, struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		}{Name: "revision", Value: imageRevision})
		err := output.Validate(context.Background())
		if err != nil {
			return input, nil, err
		}
		outputs = append(outputs, output)
	}
	return input, outputs, nil
}
func (tekton *tektonService) InitTask(step v1.Step, label map[string]string, processId string) (v1alpha1.Task, error) {
	if label == nil {
		label = make(map[string]string)
	}
	label["step"] = step.Name
	label["processId"] = processId
	task := &v1alpha1.Task{
		TypeMeta: metaV1.TypeMeta{
			Kind:       "Task",
			APIVersion: "tekton.dev/v1",
		},
		ObjectMeta: metaV1.ObjectMeta{
			Name:      step.Name + "-" + processId,
			Namespace: config.CiNamespace,
			Labels:    label,
		},
	}

	if step.Type == enums.BUILD {
		initBuildTaskSpec(step, task)
	} else if step.Type == enums.INTERMEDIARY {
		initIntermediaryTaskSpec(step, task)
	} else {
		log.Print("Please provide a valid step type!")
	}
	err := task.Validate(context.Background())
	if err != nil {
		return *task, err
	}
	return *task, nil

}
func initIntermediaryTaskSpec(step v1.Step, task *v1alpha1.Task) {
	var steps []v1alpha1.Step
	var env []corev1.EnvVar
	for key, value := range step.EnvData {
		env = append(env, corev1.EnvVar{
			Name:  key,
			Value: value,
		})
	}
	for i, image := range strings.Split(step.Params[enums.IMAGES], ",") {
		steps = append(steps, v1alpha1.Step{
			Container: corev1.Container{
				Name:            enums.CUSTOM_STAGE + strconv.Itoa(i),
				Image:           image,
				Command:         strings.Split(step.Params[enums.COMMAND], ","),
				Args:            strings.Split(step.Params[enums.COMMAND_ARGS], ","),
				Env:             env,
				ImagePullPolicy: "Always",
			}})
	}
	task.Spec.Steps = steps
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
	params = append(params, initAdditionalParams(step.ArgData)...)

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
	for i := range strings.Split(step.Params[enums.IMAGES], ",") {
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
	args := initBuildArgs(step.ArgData)
	args = append(args, "--dockerfile=$(inputs.params.pathToDockerFile)")
	args = append(args, "--context=$(inputs.params.pathToContext)")
	for i := range strings.Split(step.Params[enums.IMAGES], ",") {
		args = append(args, "--destination=$(outputs.resources.builtImage"+strconv.Itoa(i)+".url)")
		steps = append(steps, v1alpha1.Step{
			Container: corev1.Container{
				Name:    "build-and-push" + strconv.Itoa(i),
				Image:   config.KanikoImage,
				Command: []string{"/kaniko/executor"},
				Args:    args,
				Env: []corev1.EnvVar{{
					Name:  "DOCKER_CONFIG",
					Value: "/tekton/home/.docker/",
				}},
				ImagePullPolicy: "Always",
			},
		})
	}
	task.Spec.Steps = steps
}

func (tekton *tektonService) InitTaskRun(step v1.Step, label map[string]string, processId string) (v1alpha1.TaskRun, error) {
	if label == nil {
		label = make(map[string]string)
	}
	label["step"] = step.Name
	taskrun := v1alpha1.TaskRun{
		TypeMeta: metaV1.TypeMeta{
			Kind:       "TaskRun",
			APIVersion: "tekton.dev/v1alpha1",
		},
		ObjectMeta: metaV1.ObjectMeta{
			Name:      step.Name + "-" + processId,
			Namespace: config.CiNamespace,
			Labels:    label,
		},
	}
	if step.Type == enums.BUILD || step.Type == enums.INTERMEDIARY {
		taskrun.Spec = v1alpha1.TaskRunSpec{
			ServiceAccountName: step.Params[enums.SERVICE_ACCOUNT],
			TaskRef: &v1alpha1.TaskRef{
				Name: step.Name + "-" + processId,
			},
		}
	} else if step.Type == enums.JENKINS_JOB {
		taskrun.Spec = v1alpha1.TaskRunSpec{
			TaskRef: &v1alpha1.TaskRef{
				Name: enums.JENKINS_TASK_NAME,
			},
		}
	}
	if step.Type == enums.BUILD {
		var params []v1alpha1.Param
		params = append(params, v1alpha1.Param{
			Name: "pathToDockerFile",
			Value: v1alpha1.ArrayOrString{
				Type:      v1alpha1.ParamTypeString,
				StringVal: "Dockerfile",
			},
		})
		params = append(params, v1alpha1.Param{
			Name: "pathToContext",
			Value: v1alpha1.ArrayOrString{
				Type:      v1alpha1.ParamTypeString,
				StringVal: "/workspace/docker-source",
			},
		})
		if step.ArgData != nil {
			for k, v := range step.ArgData {
				params = append(params, v1alpha1.Param{
					Name: k,
					Value: v1alpha1.ArrayOrString{
						Type:      v1alpha1.ParamTypeString,
						StringVal: v,
					},
				})
			}

		}
		taskrun.Spec.Inputs = v1alpha1.TaskRunInputs{
			Resources: []v1alpha1.TaskResourceBinding{{
				PipelineResourceBinding: v1alpha1.PipelineResourceBinding{
					Name: "docker-source",
					ResourceRef: &v1alpha1.PipelineResourceRef{
						Name: step.Params[enums.REVISION][:15] + "-" + processId,
					},
				},
			}},
			Params: params,
		}
		resourceBindings := []v1alpha1.TaskResourceBinding{}
		for i := range strings.Split(step.Params[enums.IMAGES], ",") {
			resourceBindings = append(resourceBindings, v1alpha1.TaskResourceBinding{
				PipelineResourceBinding: v1alpha1.PipelineResourceBinding{
					Name: "builtImage" + strconv.Itoa(i),
					ResourceRef: &v1alpha1.PipelineResourceRef{
						Name: step.Name + "" + processId + "" + strconv.Itoa(i),
					},
				},
			})
		}
		taskRunOutputs := v1alpha1.TaskRunOutputs{
			Resources: resourceBindings,
		}
		taskrun.Spec.Outputs = taskRunOutputs
	} else if step.Type == enums.JENKINS_JOB {
		var params []v1alpha1.Param
		params = append(params, v1alpha1.Param{
			Name: "JENKINS_HOST_URL",
			Value: v1alpha1.ArrayOrString{
				Type:      v1alpha1.ParamTypeString,
				StringVal: step.Params[enums.JENKINS_URL],
			},
		})
		params = append(params, v1alpha1.Param{
			Name: "JOB_NAME",
			Value: v1alpha1.ArrayOrString{
				Type:      v1alpha1.ParamTypeString,
				StringVal: step.Params[enums.JENKINS_JOB_NAME],
			},
		})

		if step.Params[enums.JENKINS_SECRET] != "" {
			params = append(params, v1alpha1.Param{
				Name: "JENKINS_SECRETS",
				Value: v1alpha1.ArrayOrString{
					Type:      v1alpha1.ParamTypeString,
					StringVal: step.Params[enums.JENKINS_SECRET],
				},
			})
		}

		if step.Params[enums.JENKINS_PARAMS] != "" {
			paramsData := []string{}
			for _, each := range strings.Split(step.Params[enums.JENKINS_PARAMS], ",") {
				keyAndValue := strings.Split(each, ":")
				paramsData = append(paramsData, keyAndValue[0]+"="+keyAndValue[1])
			}
			params = append(params, v1alpha1.Param{
				Name: "JOB_PARAMS",
				Value: v1alpha1.ArrayOrString{
					Type:     v1alpha1.ParamTypeArray,
					ArrayVal: paramsData,
				},
			})
		}
	}
	err := taskrun.Validate(context.Background())
	if err != nil {
		return taskrun, err
	}
	return taskrun, nil
}
func (tekton *tektonService) CreatePipelineResource(resource v1alpha1.PipelineResource) error {
	_, err := tekton.Tcs.TektonV1alpha1().PipelineResources(config.CiNamespace).Create(&resource)
	if err != nil {
		log.Println("[ERROR]:", "Failed to create pipelineresource", err.Error())
		return err
	}
	return nil
}
func (tekton *tektonService) CreateTask(resource v1alpha1.Task) error {
	_, err := tekton.Tcs.TektonV1alpha1().Tasks(config.CiNamespace).Create(&resource)
	if err != nil {
		log.Println("[ERROR]:", "Failed to create task", err.Error())
		return err
	}
	return nil
}
func (tekton *tektonService) CreateTaskRun(resource v1alpha1.TaskRun) error {
	_, err := tekton.Tcs.TektonV1alpha1().TaskRuns(config.CiNamespace).Create(&resource)
	if err != nil {
		log.Println("[ERROR]:", "Failed to create taskrun", err.Error())
		return err
	}
	return nil

}
func (tekton *tektonService) DeletePipelineResourceByProcessId(processId string) error {
	list, err := tekton.Tcs.TektonV1alpha1().PipelineResources(config.CiNamespace).List(metaV1.ListOptions{
		LabelSelector: "processId=" + processId,
	})
	if err != nil {
		log.Println("[WARNING]:", err.Error())
		return err
	}
	for _, each := range list.Items {
		err = tekton.Tcs.TektonV1alpha1().PipelineResources(config.CiNamespace).Delete(each.Name, &metaV1.DeleteOptions{})
		if err != nil {
			log.Println("[ERROR]:", err.Error())
		}
	}
	return nil
}
func (tekton *tektonService) DeleteTaskByProcessId(processId string) error {
	list, err := tekton.Tcs.TektonV1alpha1().Tasks(config.CiNamespace).List(metaV1.ListOptions{
		LabelSelector: "processId=" + processId,
	})
	if err != nil {
		log.Println("[WARNING]:", err.Error())
		return err
	}
	for _, each := range list.Items {
		err = tekton.Tcs.TektonV1alpha1().Tasks(config.CiNamespace).Delete(each.Name, &metaV1.DeleteOptions{})
		if err != nil {
			log.Println("[ERROR]:", err.Error())
		}
	}
	return nil
}
func (tekton *tektonService) DeleteTaskRunByProcessId(processId string) error {
	list, err := tekton.Tcs.TektonV1alpha1().TaskRuns(config.CiNamespace).List(metaV1.ListOptions{
		LabelSelector: "processId=" + processId,
	})
	if err != nil {
		log.Println("[WARNING]:", err.Error())
		return err
	}
	for _, each := range list.Items {
		err = tekton.Tcs.TektonV1alpha1().TaskRuns(config.CiNamespace).Delete(each.Name, &metaV1.DeleteOptions{})
		if err != nil {
			log.Println("[ERROR]:", err.Error())
		}
	}
	return nil
}
func (tekton *tektonService) PurgeByProcessId(processId string) {
	_ = tekton.DeletePipelineResourceByProcessId(processId)
	_ = tekton.DeleteTaskByProcessId(processId)
	_ = tekton.DeleteTaskRunByProcessId(processId)
}
func initAdditionalParams(args map[string]string) []v1alpha1.ParamSpec {
	var params []v1alpha1.ParamSpec
	for key := range args {
		params = append(params, v1alpha1.ParamSpec{
			Name:    key,
			Type:    "string",
			Default: nil,
		})
	}
	return params
}
func initBuildArgs(arg map[string]string) []string {
	var args []string
	for key := range arg {
		args = append(args, "--build-arg="+key+"=$(inputs.params."+key+")")
	}
	return args
}

// NewTektonService returns tekton type service
func NewTektonService(tcs *versioned.Clientset) service.Tekton {
	return &tektonService{
		Tcs: tcs,
	}
}
