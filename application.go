package main

import (
	"github.com/klovercloud-ci/api"
	"github.com/klovercloud-ci/config"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

func main(){
	e:=config.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))
	api.Routes(e)
	e.Logger.Fatal(e.Start(":" + config.ServerPort))
}
