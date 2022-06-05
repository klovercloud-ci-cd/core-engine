package logic

import (
	"errors"
	v1 "github.com/klovercloud-ci-cd/core-engine/core/v1"
	"github.com/klovercloud-ci-cd/core-engine/core/v1/service"
	"github.com/klovercloud-ci-cd/core-engine/enums"
	"io"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"strconv"
)

type mockK8sService struct {
	service service.LogEvent
}

func (k8s mockK8sService) DeletePersistentVolumeClaimByProcessId(processId string) error {
	//TODO implement me
	panic("implement me")
}

func (k8s mockK8sService) CreatePersistentVolumeClaim(source corev1.PersistentVolumeClaim) error {
	//TODO implement me
	panic("implement me")
}

func (k8s mockK8sService) InitPersistentVolumeClaim(step v1.Step, label map[string]string, processId string) corev1.PersistentVolumeClaim {
	//TODO implement me
	panic("implement me")
}

func (k8s mockK8sService) RequestContainerLog(namespace string, podName string, containerName string) *rest.Request {
	panic("implement me")
}

func (k8s mockK8sService) GetContainerLog(namespace, podName, containerName string, taskRunLabel map[string]string) (io.ReadCloser, error) {
	panic("implement me")
}

func (k8s mockK8sService) FollowContainerLifeCycle(companyId, namespace, podName, containerName, step, processId string, stepType enums.STEP_TYPE, claim int) {
	panic("implement me")
}

func (k8s mockK8sService) GetPodListByProcessId(namespace, processId string, option v1.PodListGetOption) *corev1.PodList {
	panic("implement me")
}

func (k8s mockK8sService) WaitAndGetInitializedPods(companyId, namespace, processId, step, stepType string, claim int) *corev1.PodList {
	panic("implement me")
}

// GetSecret Mock k8s Secret
func (k8s *mockK8sService) GetSecret(name, namespace string) (corev1.Secret, error) {
	secrets := InitSecrets()
	for _, each := range secrets {
		if each.Name == name && each.Namespace == namespace {
			return each, nil
		}
	}
	return corev1.Secret{}, errors.New("No record found")
}

// GetConfigMap Mock k8s ConfigMap
func (k8s *mockK8sService) GetConfigMap(name, namespace string) (corev1.ConfigMap, error) {
	configMaps := InitConfigMaps()
	for _, each := range configMaps {
		if each.Name == name && each.Namespace == namespace {
			return each, nil
		}
	}
	return corev1.ConfigMap{}, errors.New("No record found")
}

// InitSecrets Mock k8s Secret list
func InitSecrets() []corev1.Secret {

	var data []corev1.Secret

	for i := 0; i < 10; i++ {
		secret := corev1.Secret{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Secret",
				APIVersion: "v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "secret" + strconv.Itoa(i),
				Namespace: "klovercloud",
			},
		}

		if i%2 == 0 {
			secret.StringData = make(map[string]string)
			secret.StringData["key1"] = "value1"
			secret.StringData["key2"] = "value2"
		} else {
			secret.Data = make(map[string][]byte)
			secret.Data["key1"] = []byte("value1")
			secret.Data["key2"] = []byte("value2")
		}
		data = append(data, secret)
	}
	return data
}

// InitConfigMaps Mock k8s ConfigMap list
func InitConfigMaps() []corev1.ConfigMap {

	var data []corev1.ConfigMap

	for i := 0; i < 10; i++ {
		configMap := corev1.ConfigMap{
			TypeMeta: metav1.TypeMeta{
				Kind:       "ConfigMap",
				APIVersion: "v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "configMap" + strconv.Itoa(i),
				Namespace: "klovercloud",
			},
			Data: map[string]string{"env1": "value1", "env2": "value2"},
		}
		data = append(data, configMap)
	}
	return data
}

// NewMockK8sService returns mock K8s type service
func NewMockK8sService(service service.LogEvent) service.K8s {
	return &mockK8sService{
		service: service,
	}
}
