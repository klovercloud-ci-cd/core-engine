package service

// Step pipeline step operations.
type Step interface {
	SetInput(url, revision string)
	SetArgs(k8s K8s)
	SetEnvs(k8s K8s)
}
