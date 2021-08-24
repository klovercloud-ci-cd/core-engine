package v1

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


type VARIABLE_REFERNECE_TYPE string
const (
	SECRET=VARIABLE_REFERNECE_TYPE("SECRET")
	CONFIG_MAP=VARIABLE_REFERNECE_TYPE("CONFIG_MAP")
)

const KLOVERCLOUD_KANIKO          = "klovercloud/kaniko:v0.14.0"