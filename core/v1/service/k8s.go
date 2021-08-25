package service
import (
	v1 "github.com/klovercloud-ci/core/v1"
	corev1 "k8s.io/api/core/v1"
)
type K8s interface {
	GetSecret(name,namespace string) (corev1.Secret,error)
	GetConfigMap(name,namespace string) (corev1.ConfigMap,error)
	GetPodListByBuildId(namespace,buildId string,option v1.PodListGetOption)*corev1.PodList
	WaitAndGetInitializedPods(namespace,buildId string)*corev1.PodList
}



