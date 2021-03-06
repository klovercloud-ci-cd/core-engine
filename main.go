package main

import (
	"github.com/klovercloud-ci-cd/core-engine/api"
	"github.com/klovercloud-ci-cd/core-engine/config"
	"github.com/klovercloud-ci-cd/core-engine/dependency"
	_ "github.com/klovercloud-ci-cd/core-engine/docs"
	"github.com/labstack/echo-contrib/jaegertracing"
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
	CollectGarbage()
	go ApplyBuildSteps()
	go ApplyIntermediarySteps()
	go ApplyJenkinsJobSteps()
	if config.EnableOpenTracing {
		c := jaegertracing.New(e, nil)
		defer c.Close()
	}
	api.Routes(e)
	e.Logger.Fatal(e.Start(":" + config.ServerPort))
}

// CollectGarbage functions that deletes garbage (tasks, taskruns, pipelineResources that ar untracked)
func CollectGarbage() {
	tektonService := dependency.GetV1TektonService()
	tektonService.DeleteTasks(config.NonPurgeAbleTasks)
	tektonService.DeleteTaskRuns(config.NonPurgeAbleTaskRuns)
	tektonService.DeletePipelineResources(config.NonPurgeAblePipelineResources)
}

// ApplyBuildSteps routine that pulls build steps in every interval.
func ApplyBuildSteps() {
	pipelineService := dependency.GetV1PipelineService()
	pipelineService.ApplyBuildSteps()
	time.Sleep(time.Second)
	ApplyBuildSteps()
}

// ApplyIntermediarySteps routine that pulls intermediary steps in every interval.
func ApplyIntermediarySteps() {
	pipelineService := dependency.GetV1PipelineService()
	pipelineService.ApplyIntermediarySteps()
	time.Sleep(time.Second)
	ApplyIntermediarySteps()
}

// ApplyJenkinsJobSteps routine that pulls intermediary steps in every interval.
func ApplyJenkinsJobSteps() {
	pipelineService := dependency.GetV1PipelineService()
	pipelineService.ApplyJenkinsJobSteps()
	time.Sleep(time.Second)
	ApplyJenkinsJobSteps()
}

// ApplyBuildCancellationSteps routine that pulls build steps in every interval.
func ApplyBuildCancellationSteps() {
	pipelineService := dependency.GetV1PipelineService()
	pipelineService.ApplyBuildSteps()
	time.Sleep(time.Second)
	ApplyBuildSteps()
}
