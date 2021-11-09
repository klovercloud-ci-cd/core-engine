package v1

import (
	"github.com/klovercloud-ci/config"
	"github.com/klovercloud-ci/dependency"
	"github.com/labstack/echo/v4"
)

// Router api/v1 base router
func Router(g *echo.Group) {
	PipelineRouter(g.Group("/pipelines"))
}

// PipelineRouter api/v1/pipeline/* router
func PipelineRouter(g *echo.Group) {
	pipelineRouter := NewPipelineApi(dependency.GetV1PipelineService(), dependency.GetV1ObserverServices())
	g.POST("", pipelineRouter.Apply, AuthenticationAndAuthorizationHandler)
	if config.UseLocalEventStore {
		g.GET("/:processId", pipelineRouter.GetLogs)
		g.GET("/ws", pipelineRouter.GetEvents)
	}

}
