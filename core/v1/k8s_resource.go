package v1

import (
	"context"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type K8s interface {
	getSecret(name,namespace string) (v1.Secret,error)
    getConfigMap(name,namespace string) (v1.ConfigMap,error)
}

type K8sResource struct {
	kcs *kubernetes.Clientset
}

func (k8s * K8sResource) getSecret(name,namespace string)(v1.Secret,error){
	sec,err:=k8s.kcs.CoreV1().
		Secrets(namespace).
		Get(context.Background(),name,metav1.GetOptions{})
	if err!=nil{
		return v1.Secret{}, err
	}
	return *sec,nil
}

func (k8s * K8sResource) getConfigMap(name,namespace string)(v1.ConfigMap,error){
	sec,err:=k8s.kcs.CoreV1().
		ConfigMaps(namespace).
		Get(context.Background(),name,metav1.GetOptions{})
	if err!=nil{
		return v1.ConfigMap{}, err
	}
	return *sec,nil
}