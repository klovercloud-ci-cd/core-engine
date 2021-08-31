package config

import (
	"github.com/joho/godotenv"
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
const (
	MONGO DATABASE= "MONGO"
	IN_MEMORY DATABASE= "IN_MEMORY"
)

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

	if CiNamespace==""{
		CiNamespace="tekton"
	}
	if KanikoImage==""{
		KanikoImage="klovercloud/kaniko:v0.14.0"
	}
	if Database==string(MONGO){
		DatabaseConnectionString = "mongodb://" + DbUsername + ":" + DbPassword + "@" + DbServer + ":" + DbPort
	}
}