package api

import (
	"github.com/labstack/echo/v4"
)

// Pipeline pipeline api operations
type Pipeline interface {
	Apply(context echo.Context) error
	GetLogs(context echo.Context) error
	GetEvents(context echo.Context) error
	CheckIfStepIsClaimable(context echo.Context) error
}
