package logic

import (
	v1 "github.com/klovercloud-ci/core/v1"
	"github.com/klovercloud-ci/core/v1/service"
	"github.com/klovercloud-ci/enums"
)

type stepService struct {
	step v1.Step
}


func (service * stepService)SetInput(url,revision string){
	if service.step.Type==enums.BUILD{
		service.step.Input.Url=url
		service.step.Input.Revision=revision
	}
}

func (service * stepService)SetArgs(k8s service.K8s){
	if service.step.Arg.Data==nil{
		service.step.Arg.Data= map[string]string{}
	}
	for _,secret:=range service.step.Arg.Secrets{
		sec,err:=k8s.GetSecret(secret.Name,secret.Namespace)
		if err!=nil{
			continue
		}
		for key,val:=range sec.StringData{
			service.step.Arg.Data[key]=string(val)
		}
	}
	for _,configMap:=range service.step.Arg.ConfigMaps{
		cm,err:=k8s.GetConfigMap(configMap.Name,configMap.Namespace)
		if err!=nil{
			continue
		}
		for key,val:=range cm.Data{
			service.step.Arg.Data[key]=val
		}

	}
}

func  (service * stepService)SetEnvs(k8s service.K8s){
	if service.step.Env.Data==nil{
		service.step.Env.Data= map[string]string{}
	}
	for _,secret:=range service.step.Env.Secrets{
		sec,err:=k8s.GetSecret(secret.Name,secret.Namespace)
		if err!=nil{
			continue
		}
		for key,val:=range sec.StringData{
			service.step.Env.Data[key]=string(val)
		}
	}
	for _,configMap:=range service.step.Env.ConfigMaps{
		cm,err:=k8s.GetConfigMap(configMap.Name,configMap.Namespace)
		if err!=nil{
			continue
		}
		for key,val:=range cm.Data{
			service.step.Env.Data[key]=val
		}

	}
}


func NewStepService(step v1.Step) service.Step {
	return &stepService{
		step: step,
	}
}
