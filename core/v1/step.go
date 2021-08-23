package v1

var _ interface {
	setInput(url,revision string)
	setArgs(k8s K8s)
	setEnvs(k8s K8s)
}= NewStep()

type Step struct {
	Name string `json:"name"`
	Type STEP_TYPE `json:"type"`
	ServiceAccount string `json:"serviceAccount"`
	Input Resource `json:"input"`
	Outputs []Resource `json:"outputs"`
	Arg Variable  `json:"arg"`
	Env Variable  `json:"env"`
}


func NewStep()*Step{
	return &Step{}
}
func (step * Step)setInput(url,revision string){
	if step.Type==BUILD{
		step.Input.Url=url
		step.Input.Revision=revision
	}
}

func (step * Step)setArgs(k8s K8s){
	if step.Arg.Data==nil{
		step.Arg.Data= map[string]string{}
	}
	for _,secret:=range step.Arg.Secrets{
		sec,err:=k8s.getSecret(secret.Name,secret.Namespace)
		if err!=nil{
			continue
		}
		for key,val:=range sec.StringData{
			step.Arg.Data[key]=string(val)
		}
	}
	for _,configMap:=range step.Arg.ConfigMaps{
		cm,err:=k8s.getConfigMap(configMap.Name,configMap.Namespace)
		if err!=nil{
			continue
		}
		for key,val:=range cm.Data{
			step.Arg.Data[key]=val
		}

	}
}

func (step * Step)setEnvs(k8s K8s){
	if step.Env.Data==nil{
		step.Env.Data= map[string]string{}
	}
	for _,secret:=range step.Env.Secrets{
		sec,err:=k8s.getSecret(secret.Name,secret.Namespace)
		if err!=nil{
			continue
		}
		for key,val:=range sec.StringData{
			step.Env.Data[key]=string(val)
		}
	}
	for _,configMap:=range step.Env.ConfigMaps{
		cm,err:=k8s.getConfigMap(configMap.Name,configMap.Namespace)
		if err!=nil{
			continue
		}
		for key,val:=range cm.Data{
			step.Env.Data[key]=val
		}

	}
}
