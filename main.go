package main

import (
	"github.com/klovercloud-ci/api"
	"github.com/klovercloud-ci/config"
	"github.com/klovercloud-ci/dependency"
	_ "github.com/klovercloud-ci/docs"
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
