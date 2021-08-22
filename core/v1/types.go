package v1

type Resource struct {
	Url string `json:"type"`
	Revision string  `json:"revision"`
}

type Variable struct {
	Secrets     []struct{
		Name string `json:"name"`
		Namespace   string `json:"namespace"`
	}`json:"secrets"`
	ConfigMaps  [] struct{
		Name string `json:"name"`
		Namespace   string `json:"namespace"`
	}`json:"configMaps"`
	Data map[string]string `json:"data"`
}