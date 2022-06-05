module github.com/klovercloud-ci-cd/core-engine

go 1.16

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/evanphx/json-patch v4.12.0+incompatible // indirect
	github.com/ghodss/yaml v1.0.0 // indirect
	github.com/go-openapi/spec v0.20.4 // indirect
	github.com/google/uuid v1.3.0
	github.com/gorilla/websocket v1.4.2
	github.com/gregjones/httpcache v0.0.0-20190611155906-901d90724c79 // indirect
	github.com/joho/godotenv v1.3.0
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/labstack/echo-contrib v0.11.0
	github.com/labstack/echo/v4 v4.6.1
	github.com/labstack/gommon v0.3.1 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mattbaird/jsonpatch v0.0.0-20171005235357-81af80346b1a // indirect
	github.com/opentracing/opentracing-go v1.2.0
	github.com/peterbourgon/diskv v2.0.1+incompatible // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/stretchr/testify v1.7.0
	github.com/swaggo/echo-swagger v1.1.4
	github.com/swaggo/files v0.0.0-20210815190702-a29dd2bc99b2 // indirect
	github.com/swaggo/swag v1.7.4
	github.com/tektoncd/pipeline v0.33.1
	go.mongodb.org/mongo-driver v1.7.1
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
	k8s.io/api v0.23.4
	k8s.io/apimachinery v0.23.4
	k8s.io/cli-runtime v0.21.4 // indirect
	k8s.io/client-go v1.5.2
	knative.dev/pkg v0.0.0-20220131144930-f4b57aef0006 // indirect
	sigs.k8s.io/yaml v1.3.0 // indirect
)

replace (
	k8s.io/api => k8s.io/api v0.22.5
	k8s.io/apimachinery => k8s.io/apimachinery v0.22.5
	k8s.io/client-go => k8s.io/client-go v0.22.5
	k8s.io/code-generator => k8s.io/code-generator v0.22.5
	k8s.io/kube-openapi => k8s.io/kube-openapi v0.0.0-20211109043538-20434351676c
)
