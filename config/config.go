package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

var IsK8 string
var CI_NAMESPACE string
var ServerPort string
var DbServer string
var DbPort string
var DbUsername string
var DbPassword string
var DatabaseConnectionString string
var DatabaseName string
var Database string


type DATABASE string
const (
	MONGO DATABASE= "MONGO"
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
	if Database==string(MONGO){
		DatabaseConnectionString = "mongodb://" + DbUsername + ":" + DbPassword + "@" + DbServer + ":" + DbPort
	}
}