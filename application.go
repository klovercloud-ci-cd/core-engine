package main

import (
	"github.com/klovercloud-ci/api"
	"github.com/klovercloud-ci/config"
	"github.com/klovercloud-ci/dependency"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"time"
)

func main(){
	e:=config.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))
	go ApplyBuildSteps()
	api.Routes(e)
	e.Logger.Fatal(e.Start(":" + config.ServerPort))
}

func ApplyBuildSteps(){
	pipelineServie:=dependency.GetPipelineService()
	pipelineServie.ApplyBuildSteps()
	time.Sleep(time.Second*5)
	ApplyBuildSteps()
}