package logic

import (
	v1 "github.com/klovercloud-ci-cd/core-engine/core/v1"
	"github.com/klovercloud-ci-cd/core-engine/core/v1/service"
	"github.com/klovercloud-ci-cd/core-engine/enums"
	"strings"
)

type stepService struct {
	step v1.Step
}

func (service *stepService) SetInput(url, revision string) {
	if service.step.Type == enums.BUILD {
		service.step.Params[enums.IMAGE_URL] = url
		service.step.Params[enums.REVISION] = revision
	}
}

func (service *stepService) SetArgs(k8s service.K8s) {
	service.step.ArgData = map[string]string{}
	if service.step.Params[enums.ARGS] != "" {
		for _, each := range strings.Split(service.step.Params[enums.ARGS], ",") {
			arg := strings.Split(each, ":")
			if len(arg) == 2 {
				service.step.ArgData[arg[0]] = arg[1]
			}
			if len(arg) == 1 {
				service.step.ArgData[arg[0]] = ""
			}
		}
	}
	for _, each := range strings.Split(service.step.Params[enums.ARGS_FROM_SECRETS], ",") {
		secret := strings.Split(each, "/")
		var name, namespace string
		if len(secret) == 2 {
			name = secret[0]
			namespace = secret[1]
		}
		if len(secret) == 1 {
			name = secret[0]
			namespace = "default"
		}
		sec, err := k8s.GetSecret(name, namespace)
		if err != nil {
			continue
		}
		for key, val := range sec.StringData {
			service.step.ArgData[key] = string(val)
		}
	}
	for _, each := range strings.Split(service.step.Params[enums.ARGS_FROM_CONFIGMAPS], ",") {
		configMap := strings.Split(each, "/")

		var name, namespace string
		if len(configMap) == 2 {
			namespace = configMap[0]
			name = configMap[1]
		}
		if len(configMap) == 1 {
			name = configMap[0]
			namespace = "default"
		}
		cm, err := k8s.GetConfigMap(name, namespace)
		if err != nil {
			continue
		}
		for key, val := range cm.Data {
			service.step.ArgData[key] = val
		}

	}
}

func (service *stepService) SetEnvs(k8s service.K8s) {
	service.step.EnvData = map[string]string{}
	if service.step.Params[enums.ENVS] != "" {
		for _, each := range strings.Split(service.step.Params[enums.ENVS], ",") {
			env := strings.Split(each, ":")
			if len(env) == 2 {
				service.step.EnvData[env[0]] = env[1]
			}
			if len(env) == 1 {
				service.step.EnvData[env[0]] = ""
			}
		}
	}

	for _, each := range strings.Split(service.step.Params[enums.ENVS_FROM_SECRETS], ",") {
		secret := strings.Split(each, "/")
		var name, namespace string
		if len(secret) == 2 {
			name = secret[0]
			namespace = secret[1]
		}
		if len(secret) == 1 {
			name = secret[0]
			namespace = "default"
		}
		sec, err := k8s.GetSecret(name, namespace)
		if err != nil {
			continue
		}
		for key, val := range sec.StringData {
			service.step.EnvData[key] = string(val)
		}
	}
	for _, each := range strings.Split(service.step.Params[enums.ENVS_FROM_CONFIGMAPS], ",") {
		configMap := strings.Split(each, "/")

		var name, namespace string
		if len(configMap) == 2 {
			namespace = configMap[0]
			name = configMap[1]
		}
		if len(configMap) == 1 {
			name = configMap[0]
			namespace = "default"
		}
		cm, err := k8s.GetConfigMap(name, namespace)
		if err != nil {
			continue
		}
		for key, val := range cm.Data {
			service.step.EnvData[key] = val
		}

	}
}

// NewStepService returns step type service.
func NewStepService(step v1.Step) service.Step {
	return &stepService{
		step: step,
	}
}
