package enums

type STEP_TYPE string


const (
	BUILD=STEP_TYPE("BUILD")
	DEPLOY=STEP_TYPE("DEPLOY")
)
const (
	Mongo = "MONGO"
	Inmemory = "INMEMORY"
)

const Api_version = "api/v1"
type PIPELINE_RESOURCE_TYPE string
const  (
	GIT=PIPELINE_RESOURCE_TYPE("git")
	IMAGE=PIPELINE_RESOURCE_TYPE("image")
)


type PIPELINE_STATUS string

const  (
	BUILD_FAILED=PIPELINE_STATUS("FAILED")
	BUILD_PROCESSING=PIPELINE_STATUS("PROCESSING")
	BUILD_TERMINATED=PIPELINE_STATUS("TERMINATED")
	WAITING=PIPELINE_STATUS("WAITING")
	TERMINATING=PIPELINE_STATUS("TERMINATING")
	INITIALIZING=PIPELINE_STATUS("INITIALIZING")
	SUCCESSFUL=PIPELINE_STATUS("SUCCESSFUL")
	CANCELLED=PIPELINE_STATUS("CANCELLED")
)

type POD_STATUS string

const (
	POD_TERMINATING      = POD_STATUS("Terminating")
	POD_INITIALIZING =POD_STATUS("PodInitializing")
)

const (
 	DEFAULT_POD_INITIALIZATION_WAIT_DURATION = 10
 	KLOVERCLOUD_KANIKO          = "klovercloud/kaniko:v0.14.0"
)




