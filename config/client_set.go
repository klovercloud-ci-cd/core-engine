package config

import (
	"flag"
	"github.com/tektoncd/pipeline/pkg/client/clientset/versioned"
	versionedResource "github.com/tektoncd/pipeline/pkg/client/resource/clientset/versioned"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
	"sync"
)

var config *rest.Config
var once sync.Once

// GetKubeConfig returns rest config for kubernetes.
func GetKubeConfig() *rest.Config {
	var config *rest.Config
	var err error
	if IsK8 == "True" {
		config, err = clientcmd.BuildConfigFromFlags("", "")
	} else {
		if home := homedir.HomeDir(); home != "" {

			kubeconfig := flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
			config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)

		} else {
			config, err = clientcmd.BuildConfigFromFlags("", "")
		}
	}
	if err != nil {
		panic(err)
	}
	return config
}

// GetClientSet returns k8s clientSets
func GetClientSet() (*versioned.Clientset,*versionedResource.Clientset, *kubernetes.Clientset,*dynamic.Interface) {
	once.Do(func() {
		config = GetKubeConfig()
	})
	cs, vcsErr := versioned.NewForConfig(config)
	vrcs,vrcsErr:=versionedResource.NewForConfig(config)
	kcs, kcsErr := kubernetes.NewForConfig(config)
	if vcsErr != nil {
		panic(vcsErr)
	}
	if kcsErr != nil {
		panic(kcsErr)
	}
	if vrcsErr!=nil{
		panic(vrcsErr)
	}
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	return cs,vrcs, kcs,&dynamicClient
}
