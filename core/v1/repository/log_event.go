package repository

import v1 "github.com/klovercloud-ci/core/v1"

type LogEventRepository interface {
	Store( log v1.LogEvent) error
	GetByBuildId(buildId string,option v1.LogEventQueryOption)[]string
}