package enums

// STEP_TYPE steps type
type STEP_TYPE string

const (
	// BUILD build step
	BUILD = STEP_TYPE("BUILD")
	// DEPLOY deploy step
	DEPLOY = STEP_TYPE("DEPLOY")
)
const (
	// MONGO mongo as db
	MONGO = "MONGO"
	// INMEMORY in memory storage as db
	INMEMORY = "INMEMORY"
)

// API_VERSION api versions
type API_VERSION string

const (
	// API_V1 api version api/v1
	API_V1 = API_VERSION("api/v1")
)

// PIPELINE_RESOURCE_TYPE pipeline resource types
type PIPELINE_RESOURCE_TYPE string

const (
	// GIT git as resource
	GIT = PIPELINE_RESOURCE_TYPE("git")
	// IMAGE docker image as resource
	IMAGE = PIPELINE_RESOURCE_TYPE("image")
	// DEPLOYMENT k8s deployment as resource
	DEPLOYMENT = PIPELINE_RESOURCE_TYPE("deployment")
	// STATEFULSET k8s statefulset as resource
	STATEFULSET = PIPELINE_RESOURCE_TYPE("statefulset")
	// DAEMONSET k8s daemonset as resource
	DAEMONSET = PIPELINE_RESOURCE_TYPE("daemonset")
	// POD k8s pod as resource
	POD = PIPELINE_RESOURCE_TYPE("pod")
	// REPLICASET k8s replicaset as resource
	REPLICASET = PIPELINE_RESOURCE_TYPE("replicaset")
)

// PIPELINE_STATUS pipeline status
type PIPELINE_STATUS string

const (
	// BUILD_FAILED step build has been FAILED
	BUILD_FAILED = PIPELINE_STATUS("FAILED")
	// BUILD_PROCESSING step build is PROCESSING
	BUILD_PROCESSING = PIPELINE_STATUS("PROCESSING")
	// BUILD_TERMINATED step build has been TERMINATED
	BUILD_TERMINATED = PIPELINE_STATUS("TERMINATED")
	// WAITING step has been in WAITING stage
	WAITING = PIPELINE_STATUS("WAITING")
	// TERMINATING step has been in TERMINATING stage
	TERMINATING = PIPELINE_STATUS("TERMINATING")
	// INITIALIZING step has been in INITIALIZING stage
	INITIALIZING = PIPELINE_STATUS("INITIALIZING")
	// SUCCESSFUL step has been SUCCESSFUL
	SUCCESSFUL = PIPELINE_STATUS("SUCCESSFUL")
	// CANCELLED step has been CANCELLED
	CANCELLED = PIPELINE_STATUS("CANCELLED")
	// ERROR step has ERROR
	ERROR = PIPELINE_STATUS("ERROR")
)

// POD_STATUS pod status
type POD_STATUS string

const (
	// POD_TERMINATING pod is Terminating
	POD_TERMINATING = POD_STATUS("Terminating")
	// POD_INITIALIZING pod is PodInitializing
	POD_INITIALIZING = POD_STATUS("PodInitializing")
)

// PIPELINE_PURGING pipeline process purging policy
type PIPELINE_PURGING string

const (
	// PIPELINE_PURGING_ENABLE pipeline process purging is enabled
	PIPELINE_PURGING_ENABLE = PIPELINE_PURGING("ENABLE")
	// PIPELINE_PURGING_DISABLE pipeline process purging is disabled
	PIPELINE_PURGING_DISABLE = PIPELINE_PURGING("DISABLE")
)

// TRIGGER pipeline trigger options
type TRIGGER string

const (
	// AUTO pipeline trigger options is auto
	AUTO = TRIGGER("AUTO")
	// MANUAL pipeline trigger options is MANUAL
	MANUAL = TRIGGER("MANUAL")
)

// PARAMS pipeline parameters
type PARAMS string

const (
	// REPOSITORY_TYPE repository type key for pipeline step param
	REPOSITORY_TYPE = PARAMS("repository_type")
	// REVISION resource revision key for  pipeline step param
	REVISION = PARAMS("revision")
	// SERVICE_ACCOUNT k8s service account key that contains registry and repository secret as pipeline step param
	SERVICE_ACCOUNT = PARAMS("service_account")
	// IMAGES key for container images as pipeline step param
	IMAGES = PARAMS("images")
	// ARGS_FROM_CONFIGMAPS key for build and other arguments via configmaps as pipeline step param
	ARGS_FROM_CONFIGMAPS = PARAMS("args_from_configmaps")
	// ARGS_FROM_SECRETS key for build and other arguments via secrets as pipeline step param
	ARGS_FROM_SECRETS = PARAMS("args_from_secrets")
	// ENVS_FROM_CONFIGMAPS key for env via configmaps as pipeline step param
	ENVS_FROM_CONFIGMAPS = PARAMS("envs_from_configmaps")
	// ENVS_FROM_SECRETS key for env via secrets as pipeline step param
	ENVS_FROM_SECRETS = PARAMS("envs_from_secrets")
	// ARGS key for build and other arguments as pipeline step param
	ARGS = PARAMS("arg")
	// ENVS key for env as pipeline step param
	ENVS = PARAMS("envs")
	// AGENT key for klovercloud-ci-agent name as pipeline step param
	AGENT = PARAMS("agent")
	// RESOURCE_NAME key for k8s resource name as pipeline step param
	RESOURCE_NAME = PARAMS("name")
	// RESOURCE_NAMESPACE key for k8s resource namespace as pipeline step param
	RESOURCE_NAMESPACE = PARAMS("namespace")
	// IMAGE_URL key for image url as pipeline step param
	IMAGE_URL = PARAMS("url")
	// TYPE key for resource type as pipeline step param
	TYPE = PARAMS("type")
)

// PROCESS_STATUS pipeline steps status
type PROCESS_STATUS string

const (
	// NON_INITIALIZED pipeline steps status non_initialized
	NON_INITIALIZED = PROCESS_STATUS("non_initialized")
	// ACTIVE pipeline steps status active
	ACTIVE = PROCESS_STATUS("active")
	// COMPLETED pipeline steps status completed
	COMPLETED = PROCESS_STATUS("completed")
	// FAILED pipeline steps status failed
	FAILED = PROCESS_STATUS("failed")
	// PAUSED pipeline steps status paused
	PAUSED = PROCESS_STATUS("paused")
)

const (
	// DEFAULT_POD_INITIALIZATION_WAIT_DURATION pod initialization wait duration for building image
	DEFAULT_POD_INITIALIZATION_WAIT_DURATION = 10
)

const (
	// DEFAULT_PAGE_LIMIT default page limit for rest api
	DEFAULT_PAGE_LIMIT = 10
	// DEFAULT_PAGE default page for rest api
	DEFAULT_PAGE = 0
)
