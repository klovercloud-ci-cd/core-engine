package v1

import (
	"context"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var _ interface {
	SetInput(url,revision string)
	SetArgs(kcs *kubernetes.Clientset)
	SetEnvs(kcs *kubernetes.Clientset)
}= NewStep()

type Step struct {
	Name string `json:"name"`
	Type STEP_TYPE `json:"type"`
	ServiceAccount string `json:"serviceAccount"`
	Input Resource `json:"input"`
	Outputs []Resource `json:"outputs"`
	Arg Variable  `json:"arg"`
	Env Variable  `json:"env"`
}


func NewStep()*Step{
	return &Step{}
}
func (step * Step)SetInput(url,revision string){
	if step.Type==BUILD{
		step.Input.Url=url
		step.Input.Revision=revision
	}
}

func (step * Step)SetArgs(kcs *kubernetes.Clientset){
	for _,secret:=range step.Arg.Secrets{
		sec,err:=getSecret(secret.Namespace,secret.Name,kcs)
		if err!=nil{
			continue
		}
		for key,val:=range sec.Data{
			step.Arg.Data[key]=string(val)
		}
	}
	for _,configMap:=range step.Arg.ConfigMaps{
		cm,err:=getConfigMap(configMap.Namespace,configMap.Name,kcs)
		if err!=nil{
			continue
		}
		for key,val:=range cm.Data{
			step.Arg.Data[key]=val
		}

	}
}

func (step * Step)SetEnvs(kcs *kubernetes.Clientset){

}

func getConfigMap(namespace,name string,kcs *kubernetes.Clientset) (v1.ConfigMap,error){
	if kcs==nil{
		//mock
	}
	cm,err:=kcs.CoreV1().
		ConfigMaps(namespace).
		Get(context.Background(),name,metav1.GetOptions{})
	if err!=nil{
		return v1.ConfigMap{}, err
	}

	return *cm,nil
}

func getSecret(namespace,name string,kcs *kubernetes.Clientset) (v1.Secret,error){
	if kcs==nil{
		//mock
	}
	sec,err:=kcs.CoreV1().
		Secrets(namespace).
		Get(context.Background(),name,metav1.GetOptions{})
	if err!=nil{
		return v1.Secret{}, err
	}

	return *sec,nil
}