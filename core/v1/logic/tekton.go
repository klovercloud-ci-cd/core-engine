package logic

import (
	"context"
	"github.com/klovercloud-ci-cd/core-engine/config"
	"github.com/klovercloud-ci-cd/core-engine/core/v1"
	"github.com/klovercloud-ci-cd/core-engine/core/v1/service"
	"github.com/klovercloud-ci-cd/core-engine/enums"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1alpha1"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	"github.com/tektoncd/pipeline/pkg/client/clientset/versioned"
	versionedResource "github.com/tektoncd/pipeline/pkg/client/resource/clientset/versioned"
	corev1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/dynamic"

	"log"
	"strconv"
	"strings"
)

type tektonService struct {
	Tcs           *versioned.Clientset
	Vrcs          *versionedResource.Clientset
	DynamicClient dynamic.Interface
	K8s           service.K8s
}

func (tekton *tektonService) GetPipelineRun(name, id, stepType string, waitUntilPipelineRunIsCompleted bool, podList corev1.PodList) (*v1beta1.PipelineRun, error) {
	pRun, pipelineRunGetingErr := tekton.Tcs.TektonV1beta1().PipelineRuns(config.CiNamespace).Get(context.Background(), name+"-"+id, metaV1.GetOptions{
		TypeMeta: metaV1.TypeMeta{
			Kind:       "PipelineRun",
			APIVersion: "tekton.dev/v1beta1",
		},
	})
	if pipelineRunGetingErr != nil {
		return nil, pipelineRunGetingErr
	}
	podListMap := make(map[string]bool)
	for _, pod := range podList.Items {
		podListMap[pod.Name] = true
	}
	if !pRun.IsDone() && waitUntilPipelineRunIsCompleted == true {
		var pList *corev1.PodList
		pList = tekton.K8s.WaitAndGetInitializedPods(config.CiNamespace, id, name)
		var NewPodList corev1.PodList
		for _, i := range pList.Items {
			if _, ok := podListMap[i.Name]; !ok {
				NewPodList.Items = append(NewPodList.Items, i)
			}
		}
		if len(NewPodList.Items) > 0 {
			for _, each := range NewPodList.Items {
				for index := range each.Spec.Containers {
					go tekton.K8s.FollowContainerLifeCycle(each.Namespace, each.Name, each.Spec.Containers[index].Name, name, id, enums.STEP_TYPE(stepType))
				}
			}

		}
		return tekton.GetPipelineRun(name, id, stepType, waitUntilPipelineRunIsCompleted, podList)
	}
	return pRun, nil
}

func (tekton *tektonService) DeletePipelineByProcessId(processId string) error {
	list, err := tekton.Tcs.TektonV1beta1().Pipelines(config.CiNamespace).List(context.Background(), metaV1.ListOptions{
		LabelSelector: "processId=" + processId,
	})
	if err != nil {
		log.Println("[WARNING]:", err.Error())
		return err
	}
	for _, each := range list.Items {
		err = tekton.Tcs.TektonV1beta1().Pipelines(config.CiNamespace).Delete(context.Background(), each.Name, metaV1.DeleteOptions{})
		if err != nil {
			log.Println("[ERROR]:", err.Error())
		}
	}
	return nil
}

func (tekton *tektonService) DeletePipelineRunByProcessId(processId string) error {
	list, err := tekton.Tcs.TektonV1beta1().PipelineRuns(config.CiNamespace).List(context.Background(), metaV1.ListOptions{
		LabelSelector: "processId=" + processId,
	})
	if err != nil {
		log.Println("[WARNING]:", err.Error())
		return err
	}
	for _, each := range list.Items {
		err = tekton.Tcs.TektonV1beta1().PipelineRuns(config.CiNamespace).Delete(context.Background(), each.Name, metaV1.DeleteOptions{})
		if err != nil {
			log.Println("[ERROR]:", err.Error())
		}
	}
	return nil
}

func (tekton *tektonService) InitPipeline(step v1.Step, label map[string]string, processId string) v1beta1.Pipeline {
	if label == nil {
		label = make(map[string]string)
	}
	label["revision"] = step.Params[enums.REVISION]
	label["processId"] = processId
	pipeline := v1beta1.Pipeline{
		TypeMeta: metaV1.TypeMeta{
			Kind:       "Pipeline",
			APIVersion: "tekton.dev/v1beta1",
		},
		ObjectMeta: metaV1.ObjectMeta{
			Name:      step.Name + "-" + processId,
			Namespace: config.CiNamespace,
			Labels:    label,
		},
		Spec: v1beta1.PipelineSpec{
			Params: []v1beta1.ParamSpec{
				{Name: "image", Type: "string", Description: "image URL to push"},
			},
			Workspaces: []v1beta1.PipelineWorkspaceDeclaration{
				{Name: "source-workspace", Description: "Directory where application source is located"},
				{Name: "cache-workspace", Description: "Directory where cache is stored for the application"},
			},
			Tasks: []v1beta1.PipelineTask{
				{Name: "fetch-repository", TaskRef: &v1beta1.TaskRef{
					Name: "git-clone",
				}, Params: []v1beta1.Param{
					{
						Name: "url", Value: v1beta1.ArrayOrString{
							Type:      v1beta1.ParamTypeString,
							StringVal: "https://github.com/buildpacks/samples",
						},
					},
					{
						Name: "subdirectory", Value: v1beta1.ArrayOrString{
							Type:      v1beta1.ParamTypeString,
							StringVal: "",
						},
					},
					{
						Name: "deleteExisting", Value: v1beta1.ArrayOrString{
							Type:      v1beta1.ParamTypeString,
							StringVal: "true",
						},
					},
				}, Workspaces: []v1beta1.WorkspacePipelineTaskBinding{{Name: "output", Workspace: "source-workspace"}}},

				{Name: "buildpacks", TaskRef: &v1beta1.TaskRef{
					Name: "buildpacks",
				}, Params: []v1beta1.Param{
					{
						Name: "APP_IMAGE", Value: v1beta1.ArrayOrString{
							Type:      v1beta1.ParamTypeString,
							StringVal: "$(params.image)",
						},
					},
					{
						Name: "SOURCE_SUBPATH", Value: v1beta1.ArrayOrString{
							Type:      v1beta1.ParamTypeString,
							StringVal: "apps/java-maven",
						},
					},
					{
						Name: "BUILDER_IMAGE", Value: v1beta1.ArrayOrString{
							Type:      v1beta1.ParamTypeString,
							StringVal: "paketobuildpacks/builder:base",
						},
					},
				}, RunAfter: []string{"fetch-repository"}, Workspaces: []v1beta1.WorkspacePipelineTaskBinding{{Name: "source", Workspace: "source-workspace"}, {Name: "cache", Workspace: "cache-workspace"}}},
			},
		},
	}
	return pipeline
}

func (tekton *tektonService) InitPipelineRun(step v1.Step, label map[string]string, processId string) (v1beta1.PipelineRun, error) {
	if label == nil {
		label = make(map[string]string)
	}
	label["revision"] = step.Params[enums.REVISION]
	label["processId"] = processId
	pipelineRun := v1beta1.PipelineRun{
		TypeMeta: metaV1.TypeMeta{
			Kind:       "PipelineRun",
			APIVersion: "tekton.dev/v1beta1",
		},
		ObjectMeta: metaV1.ObjectMeta{
			Name:      step.Name + "-" + processId,
			Namespace: config.CiNamespace,
			Labels:    label,
		},
		Spec: v1beta1.PipelineRunSpec{
			PipelineRef: &v1beta1.PipelineRef{
				Name: step.Name + "-" + processId,
			},
			Params: []v1beta1.Param{
				{
					Name: "image",
					Value: v1beta1.ArrayOrString{
						Type:      v1beta1.ParamTypeString,
						StringVal: step.Params[enums.IMAGES],
					},
				},
			},
			ServiceAccountName: step.Params[enums.SERVICE_ACCOUNT],
			Workspaces: []v1beta1.WorkspaceBinding{
				{
					Name:    "source-workspace",
					SubPath: "source",
					PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
						ClaimName: step.Name + "-" + processId,
					},
				},
				{
					Name:    "cache-workspace",
					SubPath: "cache",
					PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
						ClaimName: step.Name + "-" + processId,
					},
				},
			},
		},
	}
	err := pipelineRun.Validate(context.Background())

	if err != nil {
		return v1beta1.PipelineRun{}, err
	}
	return pipelineRun, nil
}

func (tekton *tektonService) CreatePipeline(pipeline v1beta1.Pipeline) error {
	_, err := tekton.Tcs.TektonV1beta1().Pipelines(config.CiNamespace).Create(context.Background(), &pipeline, metaV1.CreateOptions{})
	return err
}

func (tekton *tektonService) CreatePipelineRun(pipelineRun v1beta1.PipelineRun) error {
	_, err := tekton.Tcs.TektonV1beta1().PipelineRuns(config.CiNamespace).Create(context.Background(), &pipelineRun, metaV1.CreateOptions{
		TypeMeta: metaV1.TypeMeta{
			Kind:       "PipelineRun",
			APIVersion: "tekton.dev/v1beta1",
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func (tekton *tektonService) GetTaskRun(name string, waitUntilTaskRunIsCompleted bool) (*v1alpha1.TaskRun, error) {
	tRun, taskrunGetingErr := tekton.Tcs.TektonV1alpha1().TaskRuns(config.CiNamespace).Get(context.Background(), name, metaV1.GetOptions{
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
	var params []v1alpha1.ParamSpec
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
	label["processId"] = processId
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

		var paramSpecs []v1alpha1.ParamSpec
		paramSpecs = append(paramSpecs, v1alpha1.ParamSpec{
			Name: "pathToDockerFile",
			Type: v1alpha1.ParamTypeString,
			Default: &v1alpha1.ArrayOrString{
				StringVal: "Dockerfile",
			},
		})
		paramSpecs = append(paramSpecs, v1alpha1.ParamSpec{
			Name: "pathToContext",
			Type: v1alpha1.ParamTypeString,
			Default: &v1alpha1.ArrayOrString{
				StringVal: "/workspace/docker-source",
			},
		})
		if step.ArgData != nil {
			for k, v := range step.ArgData {
				paramSpecs = append(paramSpecs, v1alpha1.ParamSpec{
					Name: k,
					Type: v1alpha1.ParamTypeString,
					Default: &v1alpha1.ArrayOrString{
						StringVal: v,
					},
				})
			}
		}

		inputresourceBindings := []v1alpha1.TaskResourceBinding{}
		inputresourceBindings= append(inputresourceBindings, v1alpha1.TaskResourceBinding{
			PipelineResourceBinding: v1alpha1.PipelineResourceBinding{
				Name: "docker-source",
				ResourceRef: &v1alpha1.PipelineResourceRef{
					Name: step.Params[enums.REVISION][:15] + "-" + processId,
				},
			},
		})

		outputResourceBindings := []v1alpha1.TaskResourceBinding{}
		for i := range strings.Split(step.Params[enums.IMAGES], ",") {
			outputResourceBindings = append(outputResourceBindings, v1alpha1.TaskResourceBinding{
				PipelineResourceBinding: v1alpha1.PipelineResourceBinding{
					Name: "builtImage" + strconv.Itoa(i),
					ResourceRef: &v1alpha1.PipelineResourceRef{
						Name: step.Name + "" + processId + "" + strconv.Itoa(i),
					},
				},
			})
		}
		taskRunResources:=&v1beta1.TaskRunResources{
			Inputs:  inputresourceBindings,
			Outputs: outputResourceBindings,
		}

		taskrun.Spec.Resources=taskRunResources
		taskrun.Spec.Params=params
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
		taskrun.Spec.Params = params
		taskrun.Spec.PodTemplate.Volumes = []corev1.Volume{}
	}
	err := taskrun.Validate(context.Background())
	if err != nil {
		return taskrun, err
	}
	return taskrun, nil
}
func (tekton *tektonService) CreatePipelineResource(resource v1alpha1.PipelineResource) error {
	_, err := tekton.Vrcs.TektonV1alpha1().PipelineResources(config.CiNamespace).Create(context.Background(), &resource, metaV1.CreateOptions{})
	if err != nil {
		log.Println("[ERROR]:", "Failed to create pipelineresource", err.Error())
		return err
	}
	return nil
}
func (tekton *tektonService) CreateTask(resource v1alpha1.Task) error {
	_, err := tekton.Tcs.TektonV1alpha1().Tasks(config.CiNamespace).Create(context.Background(), &resource, metaV1.CreateOptions{})
	if err != nil {
		log.Println("[ERROR]:", "Failed to create task", err.Error())
		return err
	}
	return nil
}
func (tekton *tektonService) CreateTaskRun(resource v1alpha1.TaskRun) error {
	_, err := tekton.Tcs.TektonV1alpha1().TaskRuns(config.CiNamespace).Create(context.Background(), &resource, metaV1.CreateOptions{})
	if err != nil {
		log.Println("[ERROR]:", "Failed to create taskrun", err.Error())
		return err
	}
	return nil

}
func (tekton *tektonService) DeletePipelineResourceByProcessId(processId string) error {
	list, err := tekton.Vrcs.TektonV1alpha1().PipelineResources(config.CiNamespace).List(context.Background(), metaV1.ListOptions{
		LabelSelector: "processId=" + processId,
	})
	if err != nil {
		log.Println("[WARNING]:", err.Error())
		return err
	}
	for _, each := range list.Items {
		err = tekton.Vrcs.TektonV1alpha1().PipelineResources(config.CiNamespace).Delete(context.Background(), each.Name, metaV1.DeleteOptions{})
		if err != nil {
			log.Println("[ERROR]:", err.Error())
		}
	}
	return nil
}
func (tekton *tektonService) DeleteTaskByProcessId(processId string) error {
	list, err := tekton.Tcs.TektonV1alpha1().Tasks(config.CiNamespace).List(context.Background(), metaV1.ListOptions{
		LabelSelector: "processId=" + processId,
	})
	if err != nil {
		log.Println("[WARNING]:", err.Error())
		return err
	}
	for _, each := range list.Items {
		err = tekton.Tcs.TektonV1alpha1().Tasks(config.CiNamespace).Delete(context.Background(), each.Name, metaV1.DeleteOptions{})
		if err != nil {
			log.Println("[ERROR]:", err.Error())
		}
	}
	return nil
}
func (tekton *tektonService) DeleteTaskRunByProcessId(processId string) error {
	list, err := tekton.Tcs.TektonV1alpha1().TaskRuns(config.CiNamespace).List(context.Background(), metaV1.ListOptions{
		LabelSelector: "processId=" + processId,
	})
	if err != nil {
		log.Println("[WARNING]:", err.Error())
		return err
	}
	for _, each := range list.Items {
		err = tekton.Tcs.TektonV1alpha1().TaskRuns(config.CiNamespace).Delete(context.Background(), each.Name, metaV1.DeleteOptions{})
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
func NewTektonService(tcs *versioned.Clientset, vrcs *versionedResource.Clientset, dynamicClient *dynamic.Interface, k8s service.K8s) service.Tekton {

	return &tektonService{
		Tcs:           tcs,
		Vrcs:          vrcs,
		DynamicClient: *dynamicClient,
		K8s:           k8s,
	}
}
