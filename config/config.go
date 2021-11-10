package config

import (
	"github.com/joho/godotenv"
	"github.com/klovercloud-ci-cd/klovercloud-ci-core/enums"
	"log"
	"os"
	"strconv"
	"strings"
)

// IsK8 refers if application is running inside k8s.
var IsK8 string

// CiNamespace refers to tekton resource namespace.
var CiNamespace string

// ServerPort refers to server port.
var ServerPort string

// DbServer refers to database server ip.
var DbServer string

// DbPort refers to database server port.
var DbPort string

// DbUsername refers to database name.
var DbUsername string

// DbPassword refers to database password.
var DbPassword string

// DatabaseConnectionString refers to database connection string.
var DatabaseConnectionString string

// DatabaseName refers to database name.
var DatabaseName string

// Database refers to database options.
var Database string

// KanikoImage refers to kaniko image url.
var KanikoImage string

// EventStoreUrl refers to klovercloud-ci-event-store url.
var EventStoreUrl string

// UseLocalEventStore refers if local db will be used.
var UseLocalEventStore bool

// AllowedConcurrentBuild refers to allowed number of concurrent build.
var AllowedConcurrentBuild int64

// Publickey refers to publickey of EventStoreToken.
var Publickey string

// EnableAuthentication refers if service to service authentication is enabled.
var EnableAuthentication bool

// Token refers to jwt token for service to service communication.
var Token string

// InitEnvironmentVariables initializes environment variables
func InitEnvironmentVariables() {
	err := godotenv.Load()
	if err != nil {
		log.Println("ERROR:", err.Error())
		return
	}

	IsK8 = os.Getenv("IS_K8")
	ServerPort = os.Getenv("SERVER_PORT")
	DbServer = os.Getenv("MONGO_SERVER")
	DbPort = os.Getenv("MONGO_PORT")
	DbUsername = os.Getenv("MONGO_USERNAME")
	DbPassword = os.Getenv("MONGO_PASSWORD")
	DatabaseName = os.Getenv("DATABASE_NAME")
	Database = os.Getenv("DATABASE")
	KanikoImage = os.Getenv("KLOVERCLOUD_KANIKO")
	CiNamespace = os.Getenv("CI_NAMESPACE")
	EventStoreUrl = os.Getenv("EVENT_STORE_URL")
	if os.Getenv("USE_LOCAL_EVENT_STORE") == "true" {
		UseLocalEventStore = true
	} else {
		UseLocalEventStore = false
	}
	if CiNamespace == "" {
		CiNamespace = "tekton"
	}
	if KanikoImage == "" {
		KanikoImage = "klovercloud/kaniko:v0.14.0"
	}
	if Database == enums.MONGO {
		DatabaseConnectionString = "mongodb://" + DbUsername + ":" + DbPassword + "@" + DbServer + ":" + DbPort
	}

	AllowedConcurrentBuild, err = strconv.ParseInt(os.Getenv("ALLOWED_CONCURRENT_BUILD"), 10, 64)
	if err != nil {
		AllowedConcurrentBuild = 4
	}

	Publickey = os.Getenv("PUBLIC_KEY")

	if os.Getenv("ENABLE_AUTHENTICATION") == "" {
		EnableAuthentication = false
	} else {
		if strings.ToLower(os.Getenv("ENABLE_AUTHENTICATION")) == "true" {
			EnableAuthentication = true
		} else {
			EnableAuthentication = false
		}
	}
	Token = os.Getenv("TOKEN")
}
