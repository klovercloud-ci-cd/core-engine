package config

import (
	"github.com/joho/godotenv"
	v1 "github.com/klovercloud-ci/core/v1"
	"github.com/klovercloud-ci/enums"
	"log"
	"os"
	"strings"
)

var IsK8 string
var CiNamespace string
var ServerPort string
var DbServer string
var DbPort string
var DbUsername string
var DbPassword string
var DatabaseConnectionString string
var DatabaseName string
var Database string
var KanikoImage string
type DATABASE string
var AGENT map[string]v1.Agent
var EventStoreUrl string
var UseLocalEventStore bool

func InitEnvironmentVariables(){
	err := godotenv.Load()
	if err != nil {
		log.Println("ERROR:", err.Error())
		return
	}
	ServerPort = os.Getenv("SERVER_PORT")
	DbServer = os.Getenv("MONGO_SERVER")
	DbPort = os.Getenv("MONGO_PORT")
	DbUsername = os.Getenv("MONGO_USERNAME")
	DbPassword = os.Getenv("MONGO_PASSWORD")
	DatabaseName = os.Getenv("DATABASE_NAME")
	Database=os.Getenv("DATABASE")
	KanikoImage=os.Getenv("KLOVERCLOUD_KANIKO")
	CiNamespace =os.Getenv("CI_NAMESPACE")
	EventStoreUrl=os.Getenv("EVENT_STORE_URL")
	if os.Getenv("USE_LOCAL_EVENT_STORE")=="true"{
		UseLocalEventStore=true
	}else{
		UseLocalEventStore=false
	}
	if CiNamespace==""{
		CiNamespace="tekton"
	}
	if KanikoImage==""{
		KanikoImage="klovercloud/kaniko:v0.14.0"
	}
	if Database==enums.Mongo{
		DatabaseConnectionString = "mongodb://" + DbUsername + ":" + DbPassword + "@" + DbServer + ":" + DbPort
	}

	agents:=strings.Split(os.Getenv("AGENTS"),",")
	AGENT=make(map[string]v1.Agent)
	for _,each:=range agents{
		attrs:=strings.Split(each,"&&")
		AGENT[attrs[0]]=v1.Agent{
			Url:   attrs[1],
			Token: attrs[2],
		}
	}



}