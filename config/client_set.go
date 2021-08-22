package config

import (
	"flag"
	"github.com/tektoncd/pipeline/pkg/client/clientset/versioned"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"log"
	"path/filepath"
	"sync"
)

var config *rest.Config
var once sync.Once

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

func GetClientSet() (*versioned.Clientset, *kubernetes.Clientset) {
	once.Do(func() {
		config = GetKubeConfig()
	})

	cs, vcsErr := versioned.NewForConfig(config)
	kcs, kcsErr := kubernetes.NewForConfig(config)

	if vcsErr != nil {
		log.Printf("failed to create pipeline clientset: %s", vcsErr)
	}
	if kcsErr != nil {
		log.Printf("failed to create pipeline clientset: %s", kcsErr)
	}
	return cs, kcs
}
