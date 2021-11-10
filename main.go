package main

import (
	"github.com/klovercloud-ci-cd/klovercloud-ci-core/api"
	"github.com/klovercloud-ci-cd/klovercloud-ci-core/config"
	"github.com/klovercloud-ci-cd/klovercloud-ci-core/dependency"
	_ "github.com/klovercloud-ci-cd/klovercloud-ci-core/docs"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"time"
)

// @title Klovercloud-ci-core API
// @description Klovercloud-ci-core API
func main() {
	e := config.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))
	go ApplyBuildSteps()
	api.Routes(e)
	e.Logger.Fatal(e.Start(":" + config.ServerPort))
}

// ApplyBuildSteps routine that pulls build steps in every interval.
func ApplyBuildSteps() {
	pipelineServie := dependency.GetV1PipelineService()
	pipelineServie.ApplyBuildSteps()
	time.Sleep(time.Second * 5)
	ApplyBuildSteps()
}

//swag init --parseDependency --parseInternal
//goreportcard-cli -v
// go get golang.org/x/tools/cmd/godoc
// godoc -http=127.0.0.1:6060
// wget -r -np -N -E -p -k http://localhost:6060/pkg/github.com/klovercloud-ci/
