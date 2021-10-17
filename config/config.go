package config

import (
	"github.com/joho/godotenv"
	"github.com/klovercloud-ci/enums"
	"log"
	"os"
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
var EventStoreUrl string
var UseLocalEventStore bool
var EventStoreToken string
var PullSize string
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
	EventStoreToken=os.Getenv("EVENT_STORE_URL_TOKEN")
	PullSize=os.Getenv("PULL_SIZE")
	if PullSize==""{
		PullSize="4"
	}
}