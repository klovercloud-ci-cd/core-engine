package v1

import (
	"github.com/klovercloud-ci-cd/core-engine/config"
	"github.com/klovercloud-ci-cd/core-engine/dependency"
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
	g.GET("/:processId/steps",pipelineRouter.CheckIfStepIsClaimable)
	if config.UseLocalEventStore {
		g.GET("/:processId", pipelineRouter.GetLogs)
		g.GET("/ws", pipelineRouter.GetEvents)
	}

}
