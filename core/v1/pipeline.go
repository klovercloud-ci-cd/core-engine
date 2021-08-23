package v1

var _ interface {
	loadArgs(k8s K8s)
	loadEnvs(k8s K8s)
	setInputResource(url,revision string)
	Build(k8s K8s,url,revision string)
}= NewPipeline()


type Pipeline struct {
	ApiVersion string `json:"api_version"`
	Name string `json:"name"`
	Label map[string]string  `json:"label"`
	Steps [] Step  `json:"steps"`
}

func(pipeline *Pipeline)loadArgs(k8s K8s){
	for _,each:=range pipeline.Steps{
		each.SetArgs(k8s)
	}

}

func(pipeline *Pipeline)loadEnvs(k8s K8s){
	for _,each:=range pipeline.Steps{
		each.SetEnvs(k8s)
	}
}

func(pipeline *Pipeline)setInputResource(url,revision string){
	for _,each:=range pipeline.Steps{
		each.SetInput(url,revision)
	}
}

func(pipeline *Pipeline)Build(k8s K8s,url,revision string){
	pipeline.loadArgs(k8s)
	pipeline.loadEnvs(k8s)
	pipeline.setInputResource(url,revision)
}
func NewPipeline() *Pipeline {
	return &Pipeline{}
}
