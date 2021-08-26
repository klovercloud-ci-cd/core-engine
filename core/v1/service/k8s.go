package service
import (
	v1 "github.com/klovercloud-ci/core/v1"
	"github.com/klovercloud-ci/enums"
	corev1 "k8s.io/api/core/v1"
)
type K8s interface {
	GetSecret(name,namespace string) (corev1.Secret,error)
	GetConfigMap(name,namespace string) (corev1.ConfigMap,error)
	GetPodListByProcessId(namespace,processId string,option v1.PodListGetOption)*corev1.PodList
	WaitAndGetInitializedPods(namespace,processId,step string)*corev1.PodList
	LogContainer(namespace,podName,containerName,step,processId string,stepType enums.STEP_TYPE )
}



