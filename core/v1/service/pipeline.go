package service

type Pipeline interface {
	Apply(url,revision string)error
	LoadArgs(k8s K8s)
	LoadEnvs(k8s K8s)
	SetInputResource(url,revision string)
	Build(k8s K8s,url,revision string)
}