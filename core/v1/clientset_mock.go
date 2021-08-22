package v1

import "k8s.io/client-go/kubernetes"



func getMockClientSet() *kubernetes.Clientset{
	return &kubernetes.Clientset{}
}