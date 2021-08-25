package logic

import (
	"errors"
	v1 "github.com/klovercloud-ci/core/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strconv"
)

type mockK8sService struct {
}

func (k8s mockK8sService) GetPodListByBuildId(namespace, buildId string, option v1.PodListGetOption) *corev1.PodList {
	panic("implement me")
}

func (k8s mockK8sService) WaitAndGetInitializedPods(namespace, buildId string) *corev1.PodList {
	panic("implement me")
}

func (k8s *mockK8sService) GetSecret(name,namespace string)(corev1.Secret,error){
	secrets:=InitSecrets()
	for _,each:=range secrets{
		if each.Name==name && each.Namespace==namespace{
			return each,nil
		}
	}
return corev1.Secret{},errors.New("No record found")
}

func (k8s *mockK8sService) GetConfigMap(name,namespace string)(corev1.ConfigMap,error){
	configMaps:=InitConfigMaps()
	for _,each:=range configMaps{
		if each.Name==name && each.Namespace==namespace{
			return each,nil
		}
	}
	return corev1.ConfigMap{},errors.New("No record found")
}
func InitSecrets()[]corev1.Secret{

	var data [] corev1.Secret

	for i:=0;i<10;i++{
		secret:= corev1.Secret{
			TypeMeta:   metav1.TypeMeta{
				Kind:       "Secret",
				APIVersion: "v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:                       "secret"+strconv.Itoa(i),
				Namespace:                  "klovercloud",
			},
		}

		if i%2==0{
			secret.StringData=make(map[string]string)
			secret.StringData["key1"]= "value1"
			secret.StringData["key2"]= "value2"
		}else{
			secret.Data=make(map[string][]byte)
			secret.Data["key1"]= []byte("value1")
			secret.Data["key2"]= []byte("value2")
		}
		data= append(data,secret)
	}
	return data
}
func InitConfigMaps()[]corev1.ConfigMap{

	var data [] corev1.ConfigMap

	for i:=0;i<10;i++{
		configMap:=corev1.ConfigMap{
			TypeMeta: metav1.TypeMeta{
				Kind:       "ConfigMap",
				APIVersion: "v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "configMap" + strconv.Itoa(i),
				Namespace: "klovercloud",
			},
			Data: map[string]string{"env1":"value1","env2":"value2"},
		}
		data= append(data, configMap)
	}
	return data
}