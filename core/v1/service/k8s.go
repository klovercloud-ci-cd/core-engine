package service

import (
	v1 "github.com/klovercloud-ci-cd/klovercloud-ci-core/core/v1"
	"github.com/klovercloud-ci-cd/klovercloud-ci-core/enums"
	"io"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/rest"
)

// K8s K8s operations.
type K8s interface {
	GetSecret(name, namespace string) (corev1.Secret, error)
	GetConfigMap(name, namespace string) (corev1.ConfigMap, error)
	GetPodListByProcessId(namespace, processId string, option v1.PodListGetOption) *corev1.PodList
	WaitAndGetInitializedPods(namespace, processId, step string) *corev1.PodList
	FollowContainerLifeCycle(namespace, podName, containerName, step, processId string, stepType enums.STEP_TYPE)
	GetContainerLog(namespace, podName, containerName string, taskRunLabel map[string]string) (io.ReadCloser, error)
	RequestContainerLog(namespace string, podName string, containerName string) *rest.Request
}
