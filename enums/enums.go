package enums

type STEP_TYPE string


const (
	BUILD=STEP_TYPE("BUILD")
	DEPLOY=STEP_TYPE("DEPLOY")
)

type PIPELINE_RESOURCE_TYPE string
const  (
	GIT=PIPELINE_RESOURCE_TYPE("git")
	IMAGE=PIPELINE_RESOURCE_TYPE("image")
)




const (
 	DEFAULT_POD_INITIALIZATION_WAIT_DURATION = 10
 	KLOVERCLOUD_KANIKO          = "klovercloud/kaniko:v0.14.0"
	TERMINATING      = "Terminating"
	POD_INITIALIZING = "PodInitializing"
)

