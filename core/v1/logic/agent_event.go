package logic

import (
	"encoding/json"
	"github.com/klovercloud-ci/config"
	v1 "github.com/klovercloud-ci/core/v1"
	"github.com/klovercloud-ci/core/v1/service"
	"github.com/klovercloud-ci/enums"
	"log"
)

type agentEventService struct {
	httpPublisher service.HttpPublisher
}

func (a agentEventService) Listen(subject v1.Subject) {
	var step v1.Step
	var deploymentResources []v1.DeploymentResource
	if subject.Step==string(enums.DEPLOY){
		for i,_:=range subject.Pipeline.Steps{
			if subject.Pipeline.Steps[i].Type==enums.DEPLOY && subject.EventData["status"]==string(enums.SUCCESSFUL){
				step=subject.Pipeline.Steps[i]
			}
		}
	}
	if len(step.Outputs)>0{
		for _,each:=range step.Outputs{
			deploymentResources= append(deploymentResources,*each.DeploymentResource )
		}
	}

	if len(deploymentResources)>0{
		for _,each:=range deploymentResources{
			agentInfo:=config.AGENT[each.Agent]
			each.ProcessId=subject.Pipeline.ProcessId
			each.Step=string(enums.DEPLOY)
			header:=make(map[string]string)
			header["token"]=agentInfo.Token
			header["Content-Type"]="application/json"
			b, err := json.Marshal(each)
			if err!=nil{
				log.Println(err.Error())
				continue
			}
			go a.httpPublisher.Post(agentInfo.Url,header,b)
		}
	}
}

func NewAgentEventService(httpPublisher service.HttpPublisher) service.AgentEvent {
	return &agentEventService{
		httpPublisher: httpPublisher,
	}
}