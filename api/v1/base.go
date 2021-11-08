package v1

import (
	"github.com/klovercloud-ci/config"
	"github.com/klovercloud-ci/dependency"
	"github.com/labstack/echo/v4"
)

func Router(g *echo.Group) {
	PipelineRouter(g.Group("/pipelines"))
}


func PipelineRouter(g *echo.Group) {
	pipelineRouter := NewPipelineApi(dependency.GetPipelineService(),dependency.GetObserverServices())
	g.POST("", pipelineRouter.Apply, AuthenticationAndAuthorizationHandler)
	if config.UseLocalEventStore {
		g.GET("/:processId", pipelineRouter.GetLogs)
		g.GET("/ws", pipelineRouter.GetEvents)
	}

}