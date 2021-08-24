package v1

var _ interface {
	loadArgs(k8s K8s)
	loadEnvs(k8s K8s)
	setInputResource(url,revision string)
	Build(k8s K8s,url,revision string)
}= NewPipeline()


type Pipeline struct {
	Option PipelineApplyOption
	ApiVersion string `json:"api_version"`
	Name string `json:"name"`
	BuildId string  `json:"build_id"`
	Label map[string]string  `json:"label"`
	Steps [] Step  `json:"steps"`
}

func(pipeline Pipeline)loadArgs(k8s K8s){
	for i,_:=range pipeline.Steps{
		pipeline.Steps[i].setArgs(k8s)
	}
}

func(pipeline *Pipeline)loadEnvs(k8s K8s){
	for i,_:=range pipeline.Steps{
		pipeline.Steps[i].setEnvs(k8s)
	}
}

func(pipeline *Pipeline)setInputResource(url,revision string){
	for i,_:=range pipeline.Steps{
		pipeline.Steps[i].setInput(url,revision)
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


func(pipeline Pipeline)Validate()error{
	for _,each:=range pipeline.Steps{
		err:=each.Validate()
		if err!=nil{
			return err
		}
	}

	return nil
}