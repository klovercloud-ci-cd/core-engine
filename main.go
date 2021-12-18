package main

import (
	"github.com/klovercloud-ci-cd/core-engine/api"
	"github.com/klovercloud-ci-cd/core-engine/config"
	"github.com/klovercloud-ci-cd/core-engine/dependency"
	_ "github.com/klovercloud-ci-cd/core-engine/docs"
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
	go ApplyIntermediarySteps()
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

// ApplyIntermediarySteps routine that pulls intermediary steps in every interval.
func ApplyIntermediarySteps() {
	pipelineServie := dependency.GetV1PipelineService()
	pipelineServie.ApplyIntermediarySteps()
	time.Sleep(time.Second * 5)
	ApplyIntermediarySteps()
}


// ApplyBuildCancellationSteps routine that pulls build steps in every interval.
func ApplyBuildCancellationSteps() {
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
// goplantuml -recursive . > ClassDiagram.puml
