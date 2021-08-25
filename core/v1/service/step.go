package service
type Step interface {
	SetInput(url,revision string)
	SetArgs(k8s K8s)
	SetEnvs(k8s K8s)
}