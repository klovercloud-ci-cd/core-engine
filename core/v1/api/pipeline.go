package api

import (

	"github.com/labstack/echo/v4"
)

type Pipeline interface {
	Apply(context echo.Context) error
	GetLogs(context echo.Context) error
	GetEvents(context echo.Context) error
}
