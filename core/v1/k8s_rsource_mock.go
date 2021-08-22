package v1

import (
	"errors"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"strconv"
)

type MockK8sResource struct {
	kcs *kubernetes.Clientset
}

func (k8s * MockK8sResource) getSecret(name,namespace string)(v1.Secret,error){
	secrets:=InitSecrets()
	for _,each:=range secrets{
		if each.Name==name && each.Namespace==namespace{
			return each,nil
		}
	}
return v1.Secret{},errors.New("No record found")
}

func (k8s * MockK8sResource) getConfigMap(name,namespace string)(v1.ConfigMap,error){
	configMaps:=InitConfigMaps()
	for _,each:=range configMaps{
		if each.Name==name && each.Namespace==namespace{
			return each,nil
		}
	}
	return v1.ConfigMap{},errors.New("No record found")
}
func InitSecrets()[]v1.Secret{

	var data [] v1.Secret

	for i:=0;i<10;i++{
		secret:= v1.Secret{
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
func InitConfigMaps()[]v1.ConfigMap{

	var data [] v1.ConfigMap

	for i:=0;i<10;i++{
		configMap:=v1.ConfigMap{
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